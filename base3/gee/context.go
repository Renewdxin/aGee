package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
我们将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go
方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
*/

// H 类似于gin.H{}
type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: request,
		Path:    request.URL.Path,
		Method:  request.Method,
	}
}

//访问Query和PostForm参数的方法。

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//快速构造String/Data/JSON/HTML响应的方法。

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
