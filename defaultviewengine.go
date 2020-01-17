package expressgo

import (
	"bytes"
	"html/template"
	"path"
)

// GoViewEngine is the middleware function generator for the GO native view engine base on html/template package
func GoViewEngine(viewdir string) ViewEngine {
	vd := viewdir
	return func(templateFile string, data ViewData, resp *Response) {
		writer := bytes.NewBufferString("")
		t, err := template.ParseFiles(path.Join(vd, templateFile))
		if err != nil {
			panic(err)
		}
		err = t.Execute(writer, data)

		resp.End(writer.String())

	}
}
