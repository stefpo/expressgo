package main

import (
	"fmt"
	"math"
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
		Use(express.JSON()).
		Use("/simplePage", simpleHeader).
		Use("/view", viewtest).
		Use("/Routed", simpleHeader).
		Use("/Routed", express.Router(express.OptionsMap{"CaseSensitive": true}).Get("/page1", routedPageGET)).
		Use("/", express.Router().Get("/page1", rootPageGET).
			GetPost("/simplePage", simplePage).
			Get("/google", redir).
			Get("/json", jsonResponse).
			All("/eq2", eq2).
			Get("/error", internalError).
			Get("/params/:p1/:p2", showParams)).
		Use(express.Static(curdir+"/public", express.OptionsMap{"DefaultPage": "index.html"})).
		Listen("localhost:8080")
}

func showParams(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send("<h1>URL parameters</h1>")
	for k, v := range req.Params {
		resp.Send(k + ":" + econv.ToString(v) + "<br>")
	}

	for k, v := range req.Query {
		resp.Send(k + ":" + econv.ToString(v) + "<br>")
	}

	resp.End()
}

func eq2(req *express.Request, resp *express.Response, next func(...express.Error)) {
	a := econv.ToFloat64(req.Json["a"])
	b := econv.ToFloat64(req.Json["b"])
	c := econv.ToFloat64(req.Json["c"])

	delta := b*b - (4 * a * c)
	var sol1 float64
	var sol2 float64

	switch {
	case a == 0:
		break
	case delta > 0:
		{
			sol1 = (-b - math.Sqrt(delta)) / (2 * a)
			sol2 = (-b + math.Sqrt(delta)) / (2 * a)
		}
		break
	case delta == 0:
		{
			sol1 = -b / (2 * a)
			sol2 = sol1
		}
		break
	default:
		{
			sol1 = 0
			sol2 = 0
		}
	}

	// Return results (Use SetError to return an error object)
	resp.Json(map[string]interface{}{
		"delta": delta,
		"sol1":  sol1,
		"sol2":  sol2})
}

func jsonResponse(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.End(map[string]interface{}{"Value1": "v1", "Field2": "F2"})
}

func redir(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Redirect("http://www.google.com")
}

func internalError(req *express.Request, resp *express.Response, next func(...express.Error)) {
	panic("Fake crash")
}

func viewtest(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Render("testview.tpl", express.ViewData{"fld1": "Value 1"})
}

func simpleHeader(req *express.Request, resp *express.Response, next func(...express.Error)) {
	//panic("An error occured")
	session := req.Session()
	cnt := session.Get("hitcount")
	if cnt == "" {
		cnt = "0"
	}
	cnt = econv.ToString(econv.ToInt64(cnt) + 1)
	session.Set("hitcount", cnt)

	resp.Send("<h1>A Header</h1>")
	resp.Send("<p>HitCount: " + cnt + " id:" + session.ID + " </p>")
}

func simplePage(req *express.Request, resp *express.Response, next func(...express.Error)) {
	resp.Send("<p>Hello world ! </p>")
	if req.Method() == "POST" {
		for k, v := range req.Form {
			resp.Send(fmt.Sprintf("<b>%s</b>: %s<br>", k, v[0]))
		}
	}
	for k, v := range req.Query {
		resp.Send(fmt.Sprintf("<b>%s</b>: %s<br>", k, v))
	}
	resp.End()

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
