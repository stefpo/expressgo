package expressgo

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Error struct {
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
type ViewEngine func(templateFile string, data ViewData, resp *Response)

// Middleware structure routing and handler information of a middleware used.
type Middleware struct {
	Path    string
	Handler interface{}
}

// Application the is main application description object
type Application struct {
	Name         string
	middleware   []Middleware
	vars         map[string]interface{}
	routes       map[string]Middleware
	XPoweredBy   string
	ErrorHandler Middleware
}

// Express creates a new instance of an application
func Express() *Application {
	return &Application{
		Name:         "Basic application",
		middleware:   nil,
		vars:         make(map[string]interface{}),
		XPoweredBy:   "ExpressGo application server",
		ErrorHandler: Middleware{Path: "", Handler: defaultErrorPage}}
}

func (thisApp *Application) Set(key string, value interface{}) *Application {
	thisApp.vars[key] = value
	return thisApp
}

func (thisApp *Application) Get(key string) interface{} {
	if v, ok := thisApp.vars[key]; ok {
		return v
	} else {
		return nil
	}
}

// Use adds middleware to the application stack
func (thisApp *Application) Use(p ...interface{}) *Application {
	if len(p) == 1 {
		switch p[0].(type) {
		case (func(*Request, *Response, func(...Error))):
			mw := p[0].(func(*Request, *Response, func(...Error)))
			if thisApp.middleware == nil {
				thisApp.middleware = make([]Middleware, 0)
			}
			thisApp.middleware = append(thisApp.middleware, Middleware{Path: "", Handler: mw})
		default:
			panic("Use: Invalid type for P1 in Use. Expected func(*Request, *Response) Status")

		}
	} else if len(p) == 2 {
		switch p[0].(type) {
		case string:
			switch p[1].(type) {
			case (func(*Request, *Response, func(...Error))):
				mw := p[1].(func(*Request, *Response, func(...Error)))
				if thisApp.middleware == nil {
					thisApp.middleware = make([]Middleware, 0)
				}
				thisApp.middleware = append(thisApp.middleware, Middleware{Path: p[0].(string), Handler: mw})
			case *RouterT:
				mw := p[1].(*RouterT)
				if thisApp.middleware == nil {
					thisApp.middleware = make([]Middleware, 0)
				}
				thisApp.middleware = append(thisApp.middleware, Middleware{Path: p[0].(string), Handler: mw})
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
	App *Application
}

// ServeHTTP is the web server main handler function.
func (hdlr mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{
		Request: r,
		vars:    make(map[string]interface{}),
		App:     hdlr.App,
		session: nil,
		Json:    nil}
	resp := Response{
		writer:      w,
		ContentType: "text/html; charset=\"utf-8\"",
		status:      Error{StatusCode: http.StatusOK, Details: ""},
		App:         hdlr.App}

	resp.SetHeader("X-Powered-By", hdlr.App.XPoweredBy)

	defer func() {
		if e := recover(); e != nil {
			em := "Unhandled exception"
			switch e.(type) {
			case error:
				em = e.(error).Error()
			case string:
				em = e.(string)
			}

			if resp.status.StatusCode == 200 {
				resp.status = Error{StatusCode: http.StatusInternalServerError, Details: em}
			}
		}

		if resp.status.StatusCode != 200 {
			hdlr.App.ErrorHandler.Handler.(func(Error, *Request, *Response, func(...Error)))(resp.status, &req, &resp, func(...Error) {})
		}
	}()

	for i := 0; i < len(hdlr.App.middleware); i++ {
		middlewareIsEndpoint := true
		middlewareCalled := false

		next := func(p ...Error) {
			middlewareIsEndpoint = false
			switch len(p) {
			case 0:
				break
			case 1:
				if resp.status.StatusCode == 200 {
					resp.status = p[0]
				}
				break
			default:
				panic("next(): extra parameter")
			}
		}

		switch hdlr.App.middleware[i].Handler.(type) {
		case func(*Request, *Response, func(...Error)):
			if hdlr.App.middleware[i].Path == "" ||
				hdlr.App.middleware[i].Path == req.Path() ||
				strings.HasPrefix(req.Path(), hdlr.App.middleware[i].Path+"/") {
				middlewareCalled = true
				hdlr.App.middleware[i].Handler.(func(*Request, *Response, func(...Error)))(&req, &resp, next)
			}
			break

		case *RouterT:
			req.mountPath = hdlr.App.middleware[i].Path
			rt := hdlr.App.middleware[i].Handler.(*RouterT)
			middlewareCalled = true
			rt.handle(&req, &resp, next)
			break
		default:
			resp.status = Error{StatusCode: http.StatusNotImplemented, Details: "No handler for type"}
		}
		LogDebug(req.Path() + " " + hdlr.App.middleware[i].Path + " " + fmt.Sprintf("%d", i))
		if middlewareCalled {
			if resp.status.StatusCode != 200 || middlewareIsEndpoint {
				break
			}
		}

	}
}

// Listen starts the web server of the application
func (thisApp *Application) Listen(port string) {
	h := mainHandler{
		App: thisApp}

	if err := http.ListenAndServe(port, h); err != nil {
		fmt.Println(err.Error())
	}
}

func defaultErrorPage(err Error, req *Request, resp *Response, next func(...Error)) {
	resp.status.StatusCode = err.StatusCode
	resp.Send(fmt.Sprintf("<h1>%d %s</h1>", err.StatusCode, http.StatusText(err.StatusCode)))
	resp.Send(resp.status.Details)
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
