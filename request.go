package expressgo

import "net/http"

// Request wraps the underlying http.Request object and add a flexible data structure
// for middelware function to enrich it
type Request struct {
	*http.Request
	session   *HTTPSession
	Json      map[string]interface{}
	Params    map[string]string
	Query     map[string]string
	vars      map[string]interface{}
	App       *Application
	mountPath string
}

func (req *Request) Set(key string, value interface{}) *Request {
	req.vars[key] = value
	return req
}

func (req *Request) Get(key string) interface{} {
	if v, ok := req.vars[key]; ok {
		return v
	} else {
		return nil
	}
}

func (req *Request) Session() *HTTPSession {
	return req.session
}

// Path returns the path of the current request
func (req *Request) Path() string {
	return req.Request.URL.Path
}

// Method returns the HTTP method of the current request
func (req *Request) Method() string {
	return req.Request.Method
}

func (req *Request) UrlValue(key string, def string) string {
	if v, ok := req.Form[key]; ok {
		return v[0]
	}
	return def
}

func (req *Request) PostValue(key string, def string) string {
	_ = req.FormValue(key, "")
	if v, ok := req.PostForm[key]; ok {
		return v[0]
	}
	return def
}

func (req *Request) FormValue(key string, def string) string {
	_ = req.FormValue(key, "")
	if v, ok := req.PostForm[key]; ok {
		return v[0]
	} else {
		if v, ok := req.Form[key]; ok {
			return v[0]
		}
	}
	return def
}
