package expressgo

import (
	"log"
)

func BasicLogger(req *Request, resp *Response) Status {
	log.Printf("%s %s\n", req.Request.Method, req.Request.URL)
	return resp.OK()
}
