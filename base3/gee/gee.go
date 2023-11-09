package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 是一个处理器函数类型，用于处理 HTTP 请求。
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //中间件
	parent      *RouterGroup
	engine      *Engine //所有组共享一个实例
}

// Engine 嵌套类型
type Engine struct {
	*RouterGroup
	router *router        // 路由器，用于处理请求的函数
	groups []*RouterGroup //储存所有的路由组
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (eng *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	//要判断该请求适用于哪些中间件，在这里简单通过 URL 的前缀来判断
	for _, group := range eng.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(writer, request)
	c.handlers = middlewares
	eng.router.handle(c)
}

// 向引擎中添加路由
// - handler：处理请求的函数，将在路由匹配成功后执行
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run
// addr 是服务器监听地址。
// 返回错误对象 err。
func (eng *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, eng)
}

// Use 将中间件应用到某个group中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
