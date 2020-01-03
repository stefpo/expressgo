package expressgo

import (
	"fmt"
	"log"
	"net/http"
)

// Vars is the map structure to store request and response variables
type Vars map[string]interface{}

// HTTPStatus structure describes the HTTP response status
type HTTPStatus struct {
	StatusCode  int
	Description string
	Details     string
}

// ViewData is the data structure for passing data to a view
type ViewData map[string]string

// ViewEngine is function prototype for the view rendering function
type ViewEngine func(templateFile string, data ViewData, resp *HTTPResponse)

// HTTPRequest wraps the underlying http.Request object and add a flexible data structure
// for middelware function to enrich it
type HTTPRequest struct {
	*http.Request
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
	Vars
	HTTPStatus
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

// OK returns a 200 Success HTTP Status
func (resp *HTTPResponse) OK() HTTPStatus {
	return HTTPStatus{
		StatusCode:  200,
		Description: "Success"}
}

// Middleware structure routing and handler information of a middleware used.
type Middleware struct {
	Path    string
	Handler func(*HTTPRequest, *HTTPResponse) HTTPStatus
}

// App the is main application description object
type App struct {
	Name         string
	Middleware   []Middleware
	ViewEngine   ViewEngine
	Route        map[string]Middleware
	ErrorHandler Middleware
}

// Express creates a new instance of an application
func Express() *App {
	return &App{
		Name:         "Basic application",
		Middleware:   nil,
		ViewEngine:   nil,
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
		case (func(*HTTPRequest, *HTTPResponse) HTTPStatus):
			mw := p[0].(func(*HTTPRequest, *HTTPResponse) HTTPStatus)
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
			case (func(*HTTPRequest, *HTTPResponse) HTTPStatus):
				mw := p[1].(func(*HTTPRequest, *HTTPResponse) HTTPStatus)
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
		HTTPStatus: HTTPStatus{
			StatusCode:  200,
			Description: "OK"}}

	defer func() {
		if e := recover(); e != nil {
			em := "Unrecognized error type"
			switch e.(type) {
			case error:
				em = e.(error).Error()
			case string:
				em = e.(string)
			}

			if resp.HTTPStatus.StatusCode == 200 {
				resp.HTTPStatus = HTTPStatus{
					StatusCode:  500,
					Description: "Unhandled exception",
					Details:     em}
			}
		}

		if resp.HTTPStatus.StatusCode != 200 {
			hdlr.App.ErrorHandler.Handler(&req, &resp)
		}
	}()

	for i := 0; i < len(hdlr.App.Middleware); i++ {
		if hdlr.App.Middleware[i].Path == "" || hdlr.App.Middleware[i].Path == req.Path() {
			status := hdlr.App.Middleware[i].Handler(&req, &resp)
			if status.StatusCode != 200 {
				resp.HTTPStatus = status
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

func defaultErrorPage(req *HTTPRequest, resp *HTTPResponse) HTTPStatus {
	resp.Write(fmt.Sprintf("<h1>%d %s</h1>", resp.HTTPStatus.StatusCode, resp.HTTPStatus.Description))
	resp.Write(resp.HTTPStatus.Details)
	return resp.OK()
}

// DebugMode Sets the debugging messages on/off for the server
var DebugMode = false

// LogDebug writes debug messages if debugmode is on
func LogDebug(s string) {
	if DebugMode {
		log.Println("Debug: " + s)
	}
}
