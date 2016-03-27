package web

import (
	"fmt"
	"net/http"
)

type Context struct {
	Request   *http.Request
	Response  http.ResponseWriter
	Params    map[string]string
	webServer *Server
}

var mainServer = NewServer()
var defaultStatic []string

func Run(addr string) {
	defaultStatic = append(defaultStatic, "/static")
	mainServer.Run(addr)
}

func Get(route string, handler interface{}) {
	fmt.Println("handler= ", handler)
	mainServer.AddRoute(route, "GET", handler)
}

func POST(route string, handler interface{}) {
	mainServer.AddRoute(route, "POST", handler)
}
