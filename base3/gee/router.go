package gee

/**
实现了一个基本的路由器，它允许你注册不同的URL模式和处理函数，然后根据请求的方法和URL路径来调用相应的处理函数
我们将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go
方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
*/

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Router %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	//根据请求方法和URL路径构建一个唯一的键
	key := c.Method + "-" + c.Path
	//检查是否有与之对应的处理函数
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
