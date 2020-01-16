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
type staticServerOptions struct {
	DefaultPage string
}

func (o *staticServerOptions) merge(src map[string]interface{}) {
	setStructFromMap(o, src)
}

// Static is the middelware function generator for static file server middleware
func Static(root string, p ...OptionsMap) func(*HTTPRequest, *HTTPResponse, func(...HTTPStatus)) {
	var wwwroot = root
	options := staticServerOptions{
		DefaultPage: "index.html"}
	switch len(p) {
	case 0:
		break
	case 1:
		options.merge(p[0])
		break
	default:
		panic("Invalid arguments for Static server.")
	}
	return func(req *HTTPRequest, resp *HTTPResponse, next func(...HTTPStatus)) {
		if !resp.Complete {
			bufflen := int64(1024 * 32)

			fn := wwwroot + req.Request.URL.Path
			if req.URL.Path == "/" {
				fn = fn + options.DefaultPage
			}
			ext := filepath.Ext(fn)

			stat, err := os.Stat(fn)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				} else {
					next(HTTPStatus{StatusCode: http.StatusNotFound, Details: "File not found:" + fn})

				}
			} else {
				if stat.IsDir() {
					next(HTTPStatus{StatusCode: http.StatusForbidden, Details: "Directory listing not allowed:" + fn})
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
		next()
	}
}
