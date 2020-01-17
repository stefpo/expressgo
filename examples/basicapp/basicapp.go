package main

import (
	"fmt"
	"os"

	"github.com/stefpo/econv"
	express "github.com/stefpo/expressgo"
)

func main() {
	curdir, _ := os.Getwd()
	fmt.Println("ExpressGO basic example")
	express.DebugMode = true
	fmt.Println(curdir)

	express.Express().
		SetViewEngine(express.GoViewEngine("views")).
		Use(express.BasicLogger()).
		Use(express.Session(express.SessionConfig{Timeout: 300, CleanupInterval: 120})).
		Use(express.URLEncoded()).
		Use("/simplePage", simpleHeader).
		Use("/view", viewtest).
		Use("/Routed", simpleHeader).
		Use("/Routed", express.Router(express.OptionsMap{"CaseSensitive": true}).Get("/page1", routedPageGET)).
		Use("/", express.Router().Get("/page1", rootPageGET).
			GetPost("/simplePage", simplePage).
			Get("/google", redir).
			Get("/json", json)).
		Use(express.Static(curdir+"/public", express.OptionsMap{"DefaultPage": "index.html"})).
		Listen("localhost:8080")
}

func json(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send(23)
	resp.End(map[string]interface{}{"Value1": "v1", "Field2": "F2"})
}

func redir(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Redirect("http://www.google.com")
}

func viewtest(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Render("testview.tpl", express.ViewData{"fld1": "Value 1"})
	next()
}

func simpleHeader(req *express.Request, resp *express.Response, next func(...express.Error)) {
	//panic("An error occured")
	session, _ := req.Vars["Session"].(*express.HTTPSession)
	cnt := session.Get("hitcount")
	if cnt == "" {
		cnt = "0"
	}
	cnt = econv.ToString(econv.ToInt64(cnt) + 1)
	session.Set("hitcount", cnt)

	resp.Send("<h1>A Header</h1>")
	resp.Send("<p>HitCount: " + cnt + " id:" + session.ID + " </p>")
	next()
}

func simplePage(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send("<p>Hello world ! </p>")
	if req.Method() == "POST" {
		for k, v := range req.Form {
			resp.Send(fmt.Sprintf("<b>%s</b>: %s<br>", k, v[0]))
		}
	}
	resp.End()

	next()
}

func routedPageGET(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send("<p>I am a routed page get ! </p>")
	resp.Send("<p>URL: " + req.Path() + " </p>")
	resp.End()
}

func rootPageGET(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send("<p>I am a root page get ! </p>")
	resp.Send("<p>URL: " + req.Path() + " </p>")
	resp.End()
}
