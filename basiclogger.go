package expressgo

import (
	"log"
)

// BasicLogger is a minimal logger middleware
func BasicLogger(req *HTTPRequest, resp *HTTPResponse) HTTPStatus {
	log.Printf("%s %s\n", req.Request.Method, req.Request.URL)
	return resp.OK()
}
