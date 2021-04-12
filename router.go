package expressgo

import (
	"net/url"
	"regexp"
	"strings"
)

type routeTag struct {
	method string
	url    string
	regexp *regexp.Regexp
	params []string
}

type RouterT struct {
	routes map[*routeTag]func(*Request, *Response, func(...Error))
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
	rt.routes = make(map[*routeTag]func(*Request, *Response, func(...Error)))
	return rt
}

// Route adds a new Route for (method, url) using handler function
func (rt *RouterT) Route(method string, url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	tag := routeTag{
		method: method,
		url:    url,
		regexp: nil,
		params: []string{}}

	urlParts := strings.Split(url, "/")
	reParts := []string{}
	for level := range urlParts {
		part := urlParts[level]
		if part != "" {
			tag.params = append(tag.params, part)
			if part[0] == ':' {
				reParts = append(reParts, "(\\w+)")
			} else {
				reParts = append(reParts, "("+part+")")
			}
		}
	}

	tag.regexp = regexp.MustCompile("/" + strings.Join(reParts, "/"))
	rt.routes[&tag] = handler
	return rt
}

func (rt *RouterT) handle(req *Request, resp *Response, next func(...Error)) {
	if strings.HasPrefix(req.Path(), req.mountPath) {
		ru := strings.TrimSuffix(req.mountPath, "/")
		subPath := strings.TrimPrefix(req.Path(), ru)
		LogDebug("SubPath:" + subPath)

		httpError := Error{}
		hasError := false
		nextCalled := false

		lnext := func(p ...Error) {
			if (len(p)) > 0 {
				httpError = p[0]
				hasError = true
			}
			nextCalled = true
		}

		for k, v := range rt.routes {
			if sm := k.regexp.FindStringSubmatch(subPath); sm != nil {
				if k.regexp.MatchString(subPath) && (k.method == req.Method() || k.method == "ALL") {
					req.Params = map[string]string{}

					for i := range k.params {
						if k.params[i][0] == ':' {
							req.Params[k.params[i][1:]] = sm[i+1]
						}
					}

					req.Query = parseQueryString(req.URL.RawQuery)

					nextCalled = false
					v(req, resp, lnext)
					if !nextCalled || hasError {
						break
					}
				}
			}
		}
		if hasError {
			next(httpError)
		} else {
			next()
		}
	} else {
		next()
	}
}

func parseQueryString(query string) map[string]string {
	ret := map[string]string{}

	urlParts := strings.Split(query, "&")
	for i := range urlParts {
		p := strings.IndexByte(urlParts[i], '=')
		if p >= 0 {
			if v, e := url.QueryUnescape(urlParts[i][p+1:]); e == nil {
				ret[urlParts[i][0:p]] = v
			} else {
				ret[urlParts[i]] = ""
			}
		} else {
			ret[urlParts[i]] = ""
		}
	}

	return ret
}

// All adds new route for all methods
func (rt *RouterT) All(url string, handler func(*Request, *Response, func(...Error))) *RouterT {
	return rt.Route("ALL", url, handler)
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
