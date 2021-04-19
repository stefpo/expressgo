package main

import (
	"fmt"
	"os"

	express "github.com/stefpo/expressgo"
	"github.com/stefpo/expressgo/examples/demo/controllers"
)

func main() {
	curdir, _ := os.Getwd()
	fmt.Println("ExpressGO basic example")
	fmt.Println("Server listening on port 8080")
	express.DebugMode = true
	fmt.Println(curdir)

	express.Express().
		Set("view-engine", express.GoViewEngine("views")).
		Use(express.BasicLogger()).
		Use(express.Session(express.SessionConfig{Timeout: 300, CleanupInterval: 120})).
		Use(express.URLEncoded()).
		Use(express.JSON()).
		Use("/simplePage", controllers.SimpleHeader).
		Use("/view", controllers.Viewtest).
		Use("/Routed", controllers.SimpleHeader).
		Use("/Routed", express.Router(express.OptionsMap{"CaseSensitive": true}).Get("/page1", controllers.RoutedPageGET)).
		Use("/", express.Router().Get("/page1", controllers.RootPageGET).
			GetPost("/simplePage", controllers.SimplePage).
			Get("/google", controllers.Redir).
			Get("/json", controllers.JsonResponse).
			All("/eq2", controllers.Eq2).
			Get("/error", controllers.InternalError).
			Get("/params/:p1/:p2", controllers.ShowParams).
			Get("/html5", controllers.Html5ViewPage)).
		Use(express.Static(curdir+"/public", express.OptionsMap{"DefaultPage": "index.html"})).
		Listen("localhost:8080")
}
