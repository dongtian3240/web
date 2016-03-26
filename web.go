package web

import (
	"mime"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
)

type Context struct {
	Request *http.Request
	Params  map[string]string
	http.ResponseWriter
	Server *Server
}

func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url_ string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Header().Set("location", url_)
	ctx.ResponseWriter.Write([]byte("redirect to ", url_))
}

func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(http.StatusNotModified)
}

// not found
func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

}

func (ctx *Context) Unauthorized() {
	ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
}

func (ctx *Context) Forbidden() {
	ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
}

func (ctx *Context) ContentType(val string) string {

	var cType string
	if strings.Contains(val, "/") {
		cType = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}

		cType = mime.TypeByExtension(val)
	}

	if cType != "" {
		ctx.Header().Set("Content-Type", cType)
	}

	return cType
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.Header().Set(hdr, val)
	} else {
		ctx.Header().Add(hdr, val)
	}
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	ctx.SetHeader("Set-Cookie", cookie.String(), false)
}

var contextType reflect.Type
var defaultStaticDirs []string
var mainServer = NewServer()

var Config = &ServerConfig{
	RecoverPanic: true,
}

func init() {

	contextType = reflect.TypeOf(Context{})

	wd, _ := os.Getwd()

	defaultStaticDirs = append(defaultStaticDirs, path.Join(wd, "static"))

}

// get
func Get(route string, handler interface{}) {

	mainServer.AddRoute(route, "GET", handler)
}

// post
func Post(route string, handler interface{}) {

	mainServer.AddRoute(route, "POST", handler)
}

// put
func Put(route string, handler interface{}) {

	mainServer.AddRoute(route, "Put", handler)
}

// delete
func Delete(route string, handler interface{}) {

	mainServer.AddRoute(route, "DELETE", handler)
}
