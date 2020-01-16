package main

import (
	"fmt"
	"os"

	"github.com/stefpo/econv"
	"github.com/stefpo/expressgo"
)

func main() {
	curdir, _ := os.Getwd()
	fmt.Println("expressgo basic example")
	expressgo.DebugMode = true
	fmt.Println(curdir)

	expressgo.Express().
		SetViewEngine(expressgo.GoViewEngine("views")).
		Use(expressgo.BasicLogger).
		Use(expressgo.Session(expressgo.SessionConfig{Timeout: 300, CleanupInterval: 120})).
		Use(expressgo.URLEncoded).
		Use("/simplePage", simpleHeader).
		Use("/view", viewtest).
		Use("/Routed", simpleHeader).
		Use("/Routed", expressgo.Router(expressgo.OptionsMap{"CaseSensitive": true}).Get("/page1", routedPageGET)).
		Use("/", expressgo.Router().Get("/page1", rootPageGET).
			GetPost("/simplePage", simplePage)).
		Use(expressgo.Static(curdir+"/public", expressgo.OptionsMap{"DefaultPage": "index.html"})).
		Listen("localhost:8080")
}

func viewtest(req *expressgo.HTTPRequest, resp *expressgo.HTTPResponse, next func(...expressgo.HTTPStatus)) {
	resp.Render("testview.tpl", expressgo.ViewData{"fld1": "Value 1"})
	next()
}

func simpleHeader(req *expressgo.HTTPRequest, resp *expressgo.HTTPResponse, next func(...expressgo.HTTPStatus)) {
	//panic("An error occured")
	session, _ := req.Vars["x_session"].(*expressgo.HTTPSession)
	cnt := session.Get("hitcount")
	if cnt == "" {
		cnt = "0"
	}
	cnt = econv.ToString(econv.ToInt64(cnt) + 1)
	session.Set("hitcount", cnt)

	resp.Write("<h1>A Header</h1>")
	resp.Write("<p>HitCount: " + cnt + " id:" + session.ID + " </p>")
	next()
}

func simplePage(req *expressgo.HTTPRequest, resp *expressgo.HTTPResponse, next func(...expressgo.HTTPStatus)) {
	resp.Write("<p>Hello world ! </p>")
	if req.Method() == "POST" {
		for k, v := range req.Form {
			resp.Write(fmt.Sprintf("<b>%s</b>: %s<br>", k, v[0]))

		}
	}
	resp.End("")

	next()
}

func routedPageGET(req *expressgo.HTTPRequest, resp *expressgo.HTTPResponse, next func(...expressgo.HTTPStatus)) {
	resp.Write("<p>I am a routed page get ! </p>")
	resp.Write("<p>URL: " + req.Path() + " </p>")
	resp.End("")
}

func rootPageGET(req *expressgo.HTTPRequest, resp *expressgo.HTTPResponse, next func(...expressgo.HTTPStatus)) {
	resp.Write("<p>I am a root page get ! </p>")
	resp.Write("<p>URL: " + req.Path() + " </p>")
	resp.End("")
}
