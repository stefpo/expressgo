package expressgo

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Vars is the map structure to store request and response variables
type Vars map[string]interface{}

type HTTPStatus struct {
	StatusCode int
	Details    string
}

// ViewData is the data structure for passing data to a view
type ViewData map[string]string
type OptionsMap map[string]interface{}

func (dest OptionsMap) merge(src OptionsMap, item string) error {
	for k, v := range src {
		if x, ok := dest[k]; ok {
			if reflect.TypeOf(v) == reflect.TypeOf(x) {
				dest[k] = v
			} else {
				return fmt.Errorf("Invalid type for %s '%s'", item, k)
			}
		} else {
			return fmt.Errorf("Invalid %s '%s'", item, k)
		}
	}
	return nil
}

// ViewEngine is function prototype for the view rendering function
type ViewEngine func(templateFile string, data ViewData, resp *HTTPResponse)

// HTTPRequest wraps the underlying http.Request object and add a flexible data structure
// for middelware function to enrich it
type HTTPRequest struct {
	*http.Request
	RootURL string
	Vars
}

// Path returns the path of the current request
func (req *HTTPRequest) Path() string {
	return req.Request.URL.Path
}

// Method returns the HTTP method of the current request
func (req *HTTPRequest) Method() string {
	return req.Request.Method
}

// HTTPResponse wraps the underlying http.ResponseWriter and add a flexible data structure
// for middelware function to enrich it
// It also provides additional capabilities to manage the output
type HTTPResponse struct {
	http.ResponseWriter
	HTTPStatus
	Vars
	HeadersSent bool
	Complete    bool
	viewEngine  ViewEngine
}

// Render uses the define View engine to render data using the templateFile
func (resp *HTTPResponse) Render(templateFile string, data ViewData) {
	if resp.viewEngine != nil {
		resp.viewEngine(templateFile, data, resp)
	}
}

// End terminates the response No data will be send after that
func (resp *HTTPResponse) End(s string) {
	resp.HeadersSent = true
	resp.Write(s)
	resp.Complete = true

}

// Write sends a string to the HTTP output
func (resp *HTTPResponse) Write(s string) {
	if !resp.Complete {
		resp.HeadersSent = true
		resp.ResponseWriter.Write([]byte(s))
	}
}

// WriteBinary sends an slice of bytes to the HTTP output
func (resp *HTTPResponse) WriteBinary(d []byte) {
	if !resp.Complete {
		resp.HeadersSent = true
		resp.ResponseWriter.Write(d)
	}
}

// AddHeader adds a header to the HTTP output
func (resp *HTTPResponse) AddHeader(name string, value string) {
	resp.ResponseWriter.Header().Add(name, value)
}

// SetCookie adds a cookier to the HTTP output
func (resp *HTTPResponse) SetCookie(name string, cookie http.Cookie) {
	resp.AddHeader("Set-cookie", cookie.String())
}

// Middleware structure routing and handler information of a middleware used.
type Middleware struct {
	Path    string
	Handler interface{}
}

// App the is main application description object
type App struct {
	Name         string
	Middleware   []Middleware
	ViewEngine   ViewEngine
	Route        map[string]Middleware
	XPoweredBy   string
	ErrorHandler Middleware
}

// Express creates a new instance of an application
func Express() *App {
	return &App{
		Name:         "Basic application",
		Middleware:   nil,
		ViewEngine:   nil,
		XPoweredBy:   "ExpressGo application server",
		ErrorHandler: Middleware{Path: "", Handler: defaultErrorPage}}
}

// SetViewEngine sets the view engine for the application
func (thisApp *App) SetViewEngine(ve ViewEngine) *App {
	thisApp.ViewEngine = ve
	return thisApp
}

// Use adds middleware to the application stack
func (thisApp *App) Use(p ...interface{}) *App {
	if len(p) == 1 {
		switch p[0].(type) {
		case (func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus))):
			mw := p[0].(func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)))
			if thisApp.Middleware == nil {
				thisApp.Middleware = make([]Middleware, 0)
			}
			thisApp.Middleware = append(thisApp.Middleware, Middleware{Path: "", Handler: mw})
		default:
			panic("Use: Invalid type for P1 in Use. Expected func(*Request, *Response) Status")

		}
	} else if len(p) == 2 {
		switch p[0].(type) {
		case string:
			switch p[1].(type) {
			case (func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus))):
				mw := p[1].(func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)))
				if thisApp.Middleware == nil {
					thisApp.Middleware = make([]Middleware, 0)
				}
				thisApp.Middleware = append(thisApp.Middleware, Middleware{Path: p[0].(string), Handler: mw})
			case *RouterT:
				mw := p[1].(*RouterT)
				if thisApp.Middleware == nil {
					thisApp.Middleware = make([]Middleware, 0)
				}
				thisApp.Middleware = append(thisApp.Middleware, Middleware{Path: p[0].(string), Handler: mw})
			default:
				panic("Use: Invalid type for P2 in Use")
			}
		default:
			panic("Use: Invalid type for P1 in Use. Expected string")
		}
	}
	return thisApp
}

type mainHandler struct {
	App *App
}

// ServeHTTP is the web server main handler function.
func (hdlr mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := HTTPRequest{
		Request: r,
		Vars:    make(map[string]interface{})}
	resp := HTTPResponse{
		ResponseWriter: w,
		Vars:           make(map[string]interface{}),
		viewEngine:     hdlr.App.ViewEngine,
		HTTPStatus:     HTTPStatus{StatusCode: http.StatusOK, Details: ""}}

	resp.AddHeader("X-Powered-By", hdlr.App.XPoweredBy)

	defer func() {
		if e := recover(); e != nil {
			em := "Unhandled exception"
			switch e.(type) {
			case error:
				em = e.(error).Error()
			case string:
				em = e.(string)
			}

			if resp.HTTPStatus.StatusCode == 200 {
				resp.HTTPStatus = HTTPStatus{StatusCode: http.StatusInternalServerError, Details: em}
			}

		}

		if resp.HTTPStatus.StatusCode != 200 {
			hdlr.App.ErrorHandler.Handler.(func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)))(&req, &resp, func(...HTTPStatus) {})
		}
	}()

	for i := 0; i < len(hdlr.App.Middleware); i++ {
		middlewareIsEndpoint := true
		middlewareCalled := false

		next := func(p ...HTTPStatus) {
			middlewareIsEndpoint = false
			switch len(p) {
			case 0:
				break
			case 1:
				if resp.HTTPStatus.StatusCode == 200 {
					resp.HTTPStatus = p[0]
				}
				break
			default:
				panic("next(): extra parameter")
			}
		}

		switch hdlr.App.Middleware[i].Handler.(type) {
		case func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)):
			if hdlr.App.Middleware[i].Path == "" ||
				hdlr.App.Middleware[i].Path == req.Path() ||
				strings.HasPrefix(req.Path(), hdlr.App.Middleware[i].Path+"/") {
				middlewareCalled = true
				hdlr.App.Middleware[i].Handler.(func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)))(&req, &resp, next)
			}
			break

		case *RouterT:
			req.RootURL = hdlr.App.Middleware[i].Path
			rt := hdlr.App.Middleware[i].Handler.(*RouterT)
			middlewareCalled = true
			rt.handle(&req, &resp, next)
			break
		default:
			resp.HTTPStatus = HTTPStatus{StatusCode: http.StatusNotImplemented, Details: "No handler for type"}
		}
		LogDebug(req.Path() + " " + hdlr.App.Middleware[i].Path + " " + fmt.Sprintf("%d", i))
		if middlewareCalled {
			if resp.HTTPStatus.StatusCode != 200 || middlewareIsEndpoint {
				break
			}
		}

	}
}

// Listen starts the web server of the application
func (thisApp *App) Listen(port string) {
	h := mainHandler{
		App: thisApp}

	if err := http.ListenAndServe(port, h); err != nil {
		fmt.Println(err.Error())
	}
}

func defaultErrorPage(req *HTTPRequest, resp *HTTPResponse, next func(...HTTPStatus)) {
	resp.ResponseWriter.WriteHeader(resp.HTTPStatus.StatusCode)
	resp.Write(fmt.Sprintf("<h1>%d %s</h1>", resp.HTTPStatus.StatusCode, http.StatusText(resp.HTTPStatus.StatusCode)))
	resp.Write(resp.HTTPStatus.Details)
	next()
}

// DebugMode Sets the debugging messages on/off for the server
var DebugMode = false

// LogDebug writes debug messages if debugmode is on
func LogDebug(s string) {
	if DebugMode {
		log.Println("Debug: " + s)
	}
}
