package expressgo

import (
	"strings"
)

// URLEncoded is the middleware for parsing request body as HTML forms
func URLEncoded(req *HTTPRequest, resp *HTTPResponse) HTTPStatus {
	contentType := req.Request.Header.Get("Content-type")
	if contentType == "application/x-www-form-urlencoded" {
		if err := req.Request.ParseForm(); err != nil {
			return HTTPStatus{StatusCode: 400, Description: "Bad Request", Details: err.Error()}
		}
	}
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.Request.ParseMultipartForm(65536); err != nil {
			return HTTPStatus{StatusCode: 400, Description: "Bad Request", Details: err.Error()}
		}
	}
	return resp.OK()
}
