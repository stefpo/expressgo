package expressgo

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stefpo/econv"
)

func contentType(ext string) string {
	var mime map[string]string = map[string]string{
		".html": "text/html; charset=UTF-8",
		".htm":  "text/html; charset=UTF-8",
		".txt":  "text/plain",
		".gif":  "image/gif",
		".jpeg": "image/jpeg",
		".jpg":  "image/jpeg",
		".bmp":  "image/bmp",
		".png":  "image/png",
		".css":  "text/css",
		".json": "application/json",
		".js":   "text/javascript",
		".bin":  "application/octet_stream",
	}
	if ret, ok := mime[strings.ToLower(ext)]; ok {
		return ret
	}
	return "text/html"
}

// StaticServerConfig structure contains the static file server configuration
type StaticServerConfig struct {
	Root        string
	DefaultPage string
}

// Static is the middelware function generator for static file server middleware
func Static(config StaticServerConfig) func(req *HTTPRequest, resp *HTTPResponse) HTTPStatus {
	var wwwroot = config.Root
	var defaultPage = config.DefaultPage
	return func(req *HTTPRequest, resp *HTTPResponse) HTTPStatus {
		if !resp.Complete {
			bufflen := int64(1024 * 32)

			fn := wwwroot + req.Request.URL.Path
			if req.URL.Path == "/" {
				fn = fn + defaultPage
			}
			ext := filepath.Ext(fn)

			stat, err := os.Stat(fn)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				} else {
					return (HTTPStatus{StatusCode: http.StatusNotFound, Description: "File not found", Details: "File not found:" + fn})
				}
			} else {
				if stat.IsDir() {
					return (HTTPStatus{StatusCode: 403, Description: "Forbidden", Details: "Directory listing not allowed"})
				}
			}

			f, err := os.Open(fn)
			defer f.Close()

			size := stat.Size()
			if size < bufflen {
				bufflen = size
			}
			buff := make([]byte, bufflen)

			resp.AddHeader("Content-Length", econv.ToString(size))
			resp.AddHeader("Content-Type", contentType(ext))

			for {
				bytesRead, err := f.Read(buff)
				if err != nil {
					if err.Error() != "EOF" {
						LogDebug(err.Error())
						panic(err)
					}
				}
				if bytesRead == 0 {
					break
				}
				resp.WriteBinary(buff[:bytesRead])

			}
			resp.End("")
		}
		return resp.OK()
	}
}
