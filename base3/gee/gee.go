package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 是一个处理器函数类型，用于处理 HTTP 请求。
type HandlerFunc func(*http.Request, http.ResponseWriter)

type Engine struct {
	router map[string]HandlerFunc // 路由器，用于处理请求的函数
}

func (eng *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := eng.router[key]; ok {
		handler(request, writer)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "404 NOT FOUND: %s\n", request.URL)
	}
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 向引擎中添加路由
// - handler：处理请求的函数，将在路由匹配成功后执行
func (eng *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 将method和pattern拼接成一个key，用于存储特定的路由
	key := method + "-" + pattern
	eng.router[key] = handler
}

func (eng *Engine) GET(pattern string, handler HandlerFunc) {
	eng.addRoute("GET", pattern, handler)
}

func (eng *Engine) POST(pattern string, handler HandlerFunc) {
	eng.addRoute("POST", pattern, handler)
}

// Run
// addr 是服务器监听地址。
// 返回错误对象 err。
func (eng *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, eng)
}
