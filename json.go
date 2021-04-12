package expressgo

import (
	"encoding/json"
)

func JSON() func(req *Request, resp *Response, next func(...Error)) {
	return func(req *Request, resp *Response, next func(...Error)) {
		contentType := req.Request.Header.Get("Content-type")
		if contentType == "application/json" {
			p := make([]byte, req.ContentLength)
			_, _ = req.Body.Read(p)

			jsonData := map[string]interface{}{}

			if e := json.Unmarshal(p, &jsonData); e == nil {
				req.Json = jsonData
			} else {
				req.Json = map[string]interface{}{}
			}
		}
		next()
	}

}
