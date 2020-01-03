package expressgo

import (
	"bytes"
	"html/template"
	"path"
)

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
