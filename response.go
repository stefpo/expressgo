package expressgo

import (
	"encoding/json"
	"net/http"
)

// Response wraps the underlying http.ResponseWriter and add a flexible data structure
// for middelware function to enrich it
// It also provides additional capabilities to manage the output
type Response struct {
	App         *Application
	ContentType string
	Status      Error
	HeadersSent bool
	Vars
	isComplete bool
	viewEngine ViewEngine
	writer     http.ResponseWriter
}

// Render uses the define View engine to render data using the templateFile
func (res *Response) Render(templateFile string, data ViewData) *Response {
	if res.viewEngine != nil {
		res.viewEngine(templateFile, data, res)
	}
	return res
}

// End terminates the response No data will be send after that
func (res *Response) End(s ...interface{}) *Response {
	res.HeadersSent = true
	for i := range s {
		res.Send(s[i])
	}
	res.isComplete = true
	return res
}

// Send sends a string to the HTTP output
func (res *Response) Send(s interface{}) *Response {
	if !res.HeadersSent {
		res.writer.WriteHeader(res.Status.StatusCode)
	}
	if !res.isComplete {
		res.HeadersSent = true
		switch s.(type) {
		case string:
			res.writer.Write([]byte(s.(string)))
			break
		case []byte:
			res.writer.Write(s.([]byte))
			break
		default:
			if b, e := json.MarshalIndent(s, "", "  "); e == nil {
				res.writer.Write(b)
			} else {
				panic("Unsupported type for Response.Send()")
			}
		}
	}
	return res
}

// Set adds a header to the HTTP output
func (res *Response) Set(name string, value string) {
	res.writer.Header().Add(name, value)
}

// Cookie adds a cookier to the HTTP output
func (res *Response) Cookie(name string, cookie http.Cookie) *Response {
	res.Set("Set-cookie", cookie.String())
	return res
}

func (res *Response) Location(url string) *Response {
	res.Set("Location", url)
	return res
}

func (res *Response) Redirect(url string) *Response {
	res.Status.StatusCode = http.StatusFound
	res.Set("Refresh", "0; url="+url)
	res.End()
	return res
}
