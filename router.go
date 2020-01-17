package expressgo

import (
	"strings"
)

type routeTag struct {
	method string
	url    string
}

type RouterT struct {
	routes map[routeTag]func(*Request, *Response, func(...Error))
}

type routerOptions struct {
	CaseSensitive bool
	MergeParams   bool
	Strict        bool
}

func (o *routerOptions) merge(src map[string]interface{}) {
	setStructFromMap(o, src)
}

// Router returns an initialized instance of routerT
func Router(p ...OptionsMap) *RouterT {
	options := routerOptions{
		CaseSensitive: false,
		MergeParams:   false,
		Strict:        true}
	switch len(p) {
	case 0:
		break
	case 1:
		options.merge(p[0])
		break
	default:
		panic("Invalid arguments for Router.")
	}

	rt := new(RouterT)
	rt.routes = make(map[routeTag]func(*Request, *Response, func(...Error)))
	return rt
}

// Route adds a new route for (method, url) using handler function
func (rt *RouterT) Route(method string, url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	tag := routeTag{
		method: method,
		url:    url}
	rt.routes[tag] = handler
	return rt
}

func (rt *RouterT) handle(req *Request, resp *Response, next func(...Error)) {
	ru := strings.TrimSuffix(req.mountPath, "/")
	nextCalled := false

	lnext := func(p ...Error) {
		if (len(p)) == 0 {
			next()
		} else {
			next(p[0])
		}

		nextCalled = true
	}

	for k, v := range rt.routes {
		if ru+k.url == req.Path() && k.method == req.Method() {
			v(req, resp, lnext)
			return
		}
	}
	if !nextCalled {
		next()
	}
}

// Get adds new route for GET method
func (rt *RouterT) Get(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("GET", url, handler)
}

// Post adds new route for POST method
func (rt *RouterT) Post(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("POST", url, handler)
}

// Put adds new route for PUT method
func (rt *RouterT) Put(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("PUT", url, handler)
}

// Delete adds new route for DELETE method
func (rt *RouterT) Delete(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("DELETE", url, handler)
}

// Patch adds new route for PATCH method
func (rt *RouterT) Patch(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("PATCH", url, handler)
}

// GetPost adds new route for GET and POST methods
func (rt *RouterT) GetPost(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Get(url, handler).Post(url, handler)
}

// RESTFul adds new route for GET,POST, PUT, PATCH, DELETE methods, typcally for restful service
func (rt *RouterT) RESTFul(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Get(url, handler).
		Post(url, handler).
		Put(url, handler).
		Patch(url, handler).
		Delete(url, handler)
}
