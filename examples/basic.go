package main

import (
	"fmt"
	"os"

	"github.com/stefpo/econv"
	"github.com/stefpo/expressgo"
)

func main() {
	curdir, _ := os.Getwd()
	fmt.Println("gonnect basic example")
	app := expressgo.Express()
	expressgo.DebugMode = true
	fmt.Println(curdir)
	app.Use(expressgo.BasicLogger).
		SetViewEngine(expressgo.GoViewEngine("views")).
		Use(expressgo.Session(expressgo.SessionConfig{Timeout: 10, CleanupInterval: 60})).
		Use(expressgo.UrlEncoded).
		Use("/simplePage", simpleHeader).
		Use("/simplePage", simplePage).
		Use("/view", viewtest).
		Use(expressgo.Static(expressgo.StaticConfig{Root: curdir + "/public", DefaultPage: "index.html"}))

	app.Listen("localhost:8080")
}

func viewtest(req *expressgo.Request, resp *expressgo.Response) expressgo.Status {
	resp.Render("testview.tpl", expressgo.ViewData{"fld1": "Value 1"})
	return resp.OK()
}

func simpleHeader(req *expressgo.Request, resp *expressgo.Response) expressgo.Status {
	//panic("An error occured")
	session, _ := req.Vars["x_session"].(*expressgo.HttpSession)
	cnt := session.Get("hitcount")
	if cnt == "" {
		cnt = "0"
	}
	cnt = econv.ToString(econv.ToInt64(cnt) + 1)
	session.Set("hitcount", cnt)

	resp.Write("<h1>A Header</h1>")
	resp.Write("<p>HitCount: " + cnt + " id:" + session.Id + " </p>")
	return resp.OK()

}

func simplePage(req *expressgo.Request, resp *expressgo.Response) expressgo.Status {
	resp.Write("<p>Hello world ! </p>")
	if req.Method() == "POST" {
		for k, v := range req.Form {
			resp.Write(fmt.Sprintf("<b>%s</b>: %s<br>", k, v[0]))

		}
	}
	resp.End("")

	return resp.OK()
}
