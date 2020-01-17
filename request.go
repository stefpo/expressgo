package expressgo

import "net/http"

// Request wraps the underlying http.Request object and add a flexible data structure
// for middelware function to enrich it
type Request struct {
	*http.Request
	Vars
	App       *Application
	mountPath string
}

// Path returns the path of the current request
func (req *Request) Path() string {
	return req.Request.URL.Path
}

// Method returns the HTTP method of the current request
func (req *Request) Method() string {
	return req.Request.Method
}
