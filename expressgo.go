package expressgo

import (
	"fmt"
	"log"
	"net/http"
)

type Vars map[string]interface{}

type Status struct {
	StatusCode  int
	Description string
	Details     string
}

type ViewData map[string]string
type ViewEngine func(templateFile string, data ViewData, resp *Response)

type Request struct {
	*http.Request
	Vars
}

func (req *Request) Path() string {
	return req.Request.URL.Path
}

func (req *Request) Method() string {
	return req.Request.Method
}

type Response struct {
	http.ResponseWriter
	Vars
	Status
	HeadersSent bool
	Complete    bool
	viewEngine  ViewEngine
}

func (this *Response) Render(s string, data ViewData) {
	if this.viewEngine != nil {
		this.viewEngine(s, data, this)
	}
}

func (this *Response) End(s string) {
	this.HeadersSent = true
	this.Write(s)
	this.Complete = true

}

func (this *Response) Write(s string) {
	if !this.Complete {
		this.HeadersSent = true
		this.ResponseWriter.Write([]byte(s))
	}
}

func (this *Response) WriteBinary(d []byte) {
	if !this.Complete {
		this.HeadersSent = true
		this.ResponseWriter.Write(d)
	}
}

func (this *Response) AddHeader(name string, value string) {
	this.ResponseWriter.Header().Add(name, value)
}

func (this *Response) SetCookie(name string, cookie http.Cookie) {
	this.AddHeader("Set-cookie", cookie.String())
}

func (this *Response) OK() Status {
	return Status{
		StatusCode:  200,
		Description: "OK"}
}

type Middleware struct {
	Path    string
	Handler func(*Request, *Response) Status
}

// App the is main application object
type App struct {
	Name         string
	Middleware   []Middleware
	ViewEngine   ViewEngine
	Route        map[string]Middleware
	ErrorHandler Middleware
}

// Express create a new instance of an application
func Express() App {
	return App{
		Name:         "Basic application",
		Middleware:   nil,
		ViewEngine:   nil,
		ErrorHandler: Middleware{Path: "", Handler: defaultErrorPage}}
}

func (this *App) SetViewEngine(ve ViewEngine) *App {
	this.ViewEngine = ve
	return this
}

func (this *App) Use(p ...interface{}) *App {
	if len(p) == 1 {
		switch p[0].(type) {
		case (func(*Request, *Response) Status):
			mw := p[0].(func(*Request, *Response) Status)
			if this.Middleware == nil {
				this.Middleware = make([]Middleware, 0)
			}
			this.Middleware = append(this.Middleware, Middleware{Path: "", Handler: mw})
		default:
			panic("Use: Invalid type for P1 in Use. Expected func(*Request, *Response) Status")

		}
	} else if len(p) == 2 {
		switch p[0].(type) {
		case string:
			switch p[1].(type) {
			case (func(*Request, *Response) Status):
				mw := p[1].(func(*Request, *Response) Status)
				if this.Middleware == nil {
					this.Middleware = make([]Middleware, 0)
				}
				this.Middleware = append(this.Middleware, Middleware{Path: p[0].(string), Handler: mw})
			default:
				panic("Use: Invalid type for P2 in Use")
			}
		default:
			panic("Use: Invalid type for P1 in Use. Expected string")
		}
	}
	return this
}

func (this *App) About() string {
	return fmt.Sprintf("%d", len(this.Middleware))
}

type mainHandler struct {
	App *App
}

func (this mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{
		Request: r,
		Vars:    make(map[string]interface{})}
	resp := Response{
		ResponseWriter: w,
		Vars:           make(map[string]interface{}),
		viewEngine:     this.App.ViewEngine,
		Status: Status{
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

			if resp.Status.StatusCode == 200 {
				resp.Status = Status{
					StatusCode:  500,
					Description: "Unhandled exception",
					Details:     em}
			}
		}

		if resp.Status.StatusCode != 200 {
			this.App.ErrorHandler.Handler(&req, &resp)
		}
	}()

	for i := 0; i < len(this.App.Middleware); i++ {
		if this.App.Middleware[i].Path == "" || this.App.Middleware[i].Path == req.Path() {
			status := this.App.Middleware[i].Handler(&req, &resp)
			if status.StatusCode != 200 {
				resp.Status = status
				break
			}
		}
	}
}

func (this *App) Listen(port string) {
	h := mainHandler{
		App: this}

	if err := http.ListenAndServe(port, h); err != nil {
		fmt.Println(err.Error())
	}
}

func defaultErrorPage(req *Request, resp *Response) Status {
	resp.Write(fmt.Sprintf("<h1>%d %s</h1>", resp.Status.StatusCode, resp.Status.Description))
	resp.Write(resp.Status.Details)
	return resp.OK()
}

var DebugMode = false

func LogDebug(s string) {
	if DebugMode {
		log.Println("Debug: " + s)
	}
}
