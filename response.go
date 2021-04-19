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
	status      Error
	headersSent bool
	isComplete  bool
	viewEngine  ViewEngine
	writer      http.ResponseWriter
}

// Render uses the define View engine to render data using the templateFile
func (res *Response) Render(templateFile string, data ViewData) *Response {
	viewEngine := res.App.Get("view-engine")
	viewEngine.(ViewEngine)(templateFile, data, res)
	return res
}

// End terminates the response No data will be send after that
func (res *Response) End(s ...interface{}) {
	for i := range s {
		res.Send(s[i])
	}
	res.isComplete = true
}

func (res *Response) Json(s interface{}) {
	if b, e := json.Marshal(s); e == nil {
		res.sendBytes(b)
	} else {
		res.sendBytes([]byte("{}"))
	}

	res.isComplete = true
}

// Send sends a string to the HTTP output
func (res *Response) Send(s interface{}) *Response {
	if !res.isComplete {
		switch s.(type) {
		case string:
			res.sendBytes([]byte(s.(string)))
			break
		case []byte:
			res.sendBytes(s.([]byte))
			break
		default:
			if b, e := json.MarshalIndent(s, "", "  "); e == nil {
				res.sendBytes(b)
			} else {
				res.status.StatusCode = http.StatusInternalServerError
				res.sendBytes([]byte("Response.Send(). Unsupported data type."))
				panic("Response.Send(). Unsupported data type. " + e.Error())
			}
		}
	}
	return res
}

func (res *Response) sendBytes(b []byte) {
	if !res.headersSent {
		res.writer.WriteHeader(res.status.StatusCode)
		res.headersSent = true
	}
	res.writer.Write(b)
}

// Set adds a header to the HTTP output
func (res *Response) SetHeader(name string, value string) {
	res.writer.Header().Add(name, value)
}

// Cookie adds a cookier to the HTTP output
func (res *Response) Cookie(name string, cookie http.Cookie) *Response {
	res.SetHeader("Set-cookie", cookie.String())
	return res
}

func (res *Response) Location(url string) *Response {
	res.SetHeader("Location", url)
	return res
}

func (res *Response) Redirect(url string) *Response {
	res.status.StatusCode = http.StatusFound
	res.SetHeader("Refresh", "0; url="+url)
	res.End()
	return res
}
