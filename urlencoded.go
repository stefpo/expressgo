package expressgo

import (
	"strings"
)

func UrlEncoded(req *Request, resp *Response) Status {
	contentType := req.Request.Header.Get("Content-type")
	if contentType == "application/x-www-form-urlencoded" {
		if err := req.Request.ParseForm(); err != nil {
			return Status{StatusCode: 400, Description: "Bad Request", Details: err.Error()}
		}
	}
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.Request.ParseMultipartForm(65536); err != nil {
			return Status{StatusCode: 400, Description: "Bad Request", Details: err.Error()}
		}
	}
	return resp.OK()
}
