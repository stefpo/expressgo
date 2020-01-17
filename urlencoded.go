package expressgo

import (
	"net/http"
	"strings"
)

// URLEncoded is the middleware for parsing request body as HTML forms
func URLEncoded() func(*Request, *Response, func(...Error)) {
	return func(req *Request, resp *Response, next func(...Error)) {
		contentType := req.Request.Header.Get("Content-type")
		if contentType == "application/x-www-form-urlencoded" {
			if err := req.Request.ParseForm(); err != nil {
				next(Error{StatusCode: http.StatusBadRequest, Details: err.Error()})
			}
		}
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := req.Request.ParseMultipartForm(65536); err != nil {
				next(Error{StatusCode: http.StatusBadRequest, Details: err.Error()})
			}
		}
		next()
	}
}
