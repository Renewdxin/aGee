package gee

/**
实现了一个基本的路由器，它允许你注册不同的URL模式和处理函数，然后根据请求的方法和URL路径来调用相应的处理函数
我们将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go
方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
*/

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 用于将路由模式字符串解析为易于处理的部分，以便路由器能够根据请求的 URL 进行路由匹配
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// 在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params。这样就能够在handler中，通过Context对象访问到具体的值了。
// handle 函数中，将从路由匹配得到的 Handler 添加到 c.handlers列表中，执行c.Next()
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(context *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND %S\n", c.Path)
		})
	}
	c.Next()
}

//Test function
//func newTestRouter() *router {
//	r := newRouter()
//	r.addRoute("GET", "/", nil)
//	r.addRoute("GET", "/hello/:name", nil)
//	r.addRoute("GET", "/hello/b/c", nil)
//	r.addRoute("GET", "/hi/:name", nil)
//	r.addRoute("GET", "/assets/*filepath", nil)
//	return r
//}
//
//func TestParsePattern(t *testing.T) {
//	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
//	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
//	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
//	if !ok {
//		t.Fatal("test parsePattern failed")
//	}
//}
//
//func TestGetRoute(t *testing.T) {
//	r := newTestRouter()
//	n, ps := r.getRoute("GET", "/hello/geektutu")
//
//	if n == nil {
//		t.Fatal("nil shouldn't be returned")
//	}
//
//	if n.pattern != "/hello/:name" {
//		t.Fatal("should match /hello/:name")
//	}
//
//	if ps["name"] != "geektutu" {
//		t.Fatal("name should be equal to 'geektutu'")
//	}
//
//	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
//
//}
