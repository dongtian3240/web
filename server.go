package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type Server struct {
	Addr     string
	Port     int
	routes   []Route
	listener net.Listener
}

type Route struct {
	url         string
	method      string
	reg         *regexp.Regexp
	httpHandler http.Handler
	handler     reflect.Value
}

func NewServer() *Server {

	return &Server{
		Addr: "localhost",
		Port: 8080,
	}

}

func (s *Server) Run(addr string) {

	mux := http.NewServeMux()
	mux.Handle("/", s)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listenAndServer :", err)
	} else {
		log.Println("listener is success!")
	}
	s.listener = ls
	log.Println("start http server !")
	err = http.Serve(s.listener, mux)
	log.Println("http server is success!")
	s.listener.Close()
}

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	s.Process(response, request)
}

func (s *Server) Process(response http.ResponseWriter, request *http.Request) {

	route := s.RouteHandler(response, request)
	if route != nil {
		fmt.Println(route.httpHandler)
		route.httpHandler.ServeHTTP(response, request)
	}
}

func (s *Server) RouteHandler(response http.ResponseWriter, request *http.Request) *Route {

	requestPath := request.URL.Path
	log.Println("path", requestPath)
	ctx := &Context{Request: request, Response: response, Params: map[string]string{}, webServer: s}

	//tm := time.Now().UTC()
	//将解析的放入 r.Form 中
	request.ParseForm()

	if len(request.Form) > 0 {
		for k, v := range request.Form {

			ctx.Params[k] = v[0]
		}
	}

	//尝试是否静态资源
	if request.Method == "GET" || request.Method == "Head" {

		if s.StaticFile(requestPath, response, request) {
			return nil
		}
	}

	//查找路由
	log.Println("search router.....")
	for i := 0; i < len(s.routes); i++ {

		route := s.routes[i]

		//判断方法是否匹配
		if request.Method != route.method {
			continue
		}

		fmt.Println("======ddddd==", route)
		//正则没有匹配
		if !route.reg.MatchString(requestPath) {
			continue
		}

		match := route.reg.FindStringSubmatch(requestPath)
		//不是全匹配
		if len(match[0]) != len(requestPath) {
			continue
		}

		if route.httpHandler != nil {
			fmt.Println("匹配到=====", route)
			return &route
		}

	}
	return nil
}

func (s *Server) AddRoute(url, method string, handler interface{}) {

	re, er := regexp.Compile(url)
	if er != nil {
		return
	}
	switch handler.(type) {

	case http.Handler:
		log.Println("========匹配到http handler type ")
		s.routes = append(s.routes, Route{url: url, method: method, reg: re, httpHandler: handler.(http.Handler)})

	default:
		log.Println("match the handler.....")
		fv := reflect.ValueOf(handler)
		s.routes = append(s.routes, Route{url: url, method: method, reg: re, handler: fv})

	}

}

func (s *Server) StaticFile(name string, response http.ResponseWriter, request *http.Request) bool {

	wd, _ := os.Getwd()
	for _, staticDir := range defaultStatic {
		var indx = strings.Index(name, staticDir)
		if indx == 0 {
			staticFile := path.Join(wd, name)
			log.Println("staticFile", staticFile)
			if FileExists(staticFile) {
				log.Println("has static file")
				http.ServeFile(response, request, staticFile)
				return true
			}
		}

	}

	return false
}
