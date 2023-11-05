# The first day 11.4

## request.URL.Path 和request.URL有什么区别
request.URL是一个结构体，其中包含了请求的URL的完整信息，包括协议、主机、路径、查询参数等。
而request.URL.Path只是URL的路径部分，不包括其他信息。
所以区别就是request.URL.Path只包含URL的路径部分，而request.URL包含了完整的URL信息。

## http.ListenAndServe(address string, handler Handler) error
第二个参数类型是接口类型 http.Handler,是从 http 的源码中找到的。
Go 语言中，实现了接口方法的 struct 都可以强制转换为接口类型。
```go
type Handler interface {
ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```
## 关于go.mod 和 go.sum

### 关于go.mod 
go.mod 文件是一个文本文件，其中包含项目的依赖项和其版本，以及项目的模块定义。它由 Go 命令生成和更新，用于记录项目的依赖关系。

### 关于go.sum
go.sum 文件是一个二进制文件，其中包含依赖项的校验和，用于验证依赖项的完整性及一致性。它由 Go 命令生成和更新，不可编辑。

### 区别
这两个文件是 Go 项目中非常重要的一部分
go.mod 用于管理项目的依赖关系，而 go.sum 用于验证依赖项的完整性及一致性，确保项目依赖的正确性。


# the second day 11.5

## TASK
1. 将路由(router)独立出来，方便之后增强。
2. 设计上下文(Context)，封装 Request 和 Response ，提供对 JSON、HTML 等返回类型的支持。
3. 动手写 Gee 框架的第二天，框架代码140行，新增代码约90行
第二天成果展示
```go
func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
```

## 设计Context的原因
在构造一个完整的响应时，需要包含很多内容，如果没有有效的封装，则需要进行大量复杂的构造
1. 简化接口调用
2. 方便支撑额外功能，如解析动态路由，支持中间件
设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。 路由的处理函数，以及将要实现的中间件，参数都统一使用 Context 实例， Context 就像一次会话的百宝箱，可以找到任何东西。


## new函数为什么返回指针
1. 需要在函数内部修改对象的状态：如果你希望在函数内部修改对象的状态，并且这些更改应该在函数调用结束后保持有效，那么你通常应该返回指向对象的指针。返回对象的值将创建对象的副本，而不会影响原始对象。 
2. 避免复制开销：在Go中，将大型对象作为值传递会导致对象的复制，这可能会导致性能问题。通过返回指针，可以避免复制整个对象，从而提高性能。