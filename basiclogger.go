package expressgo

import (
	"log"
)

// BasicLogger is a minimal logger middleware
func BasicLogger() func(*Request, *Response, func(...Error)) {
	return func(req *Request, resp *Response, next func(...Error)) {
		log.Printf("%s %s\n", req.Request.Method, req.Request.URL)
		next()
	}
}
