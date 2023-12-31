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

# the third day 11.6

## TASK
1. 使用 Trie 树实现动态路由(dynamic route)解析。
2. 支持两种模式:name和*filepath，代码约150行。

## Trie 树简介

键值对的存储的方式，只能用来索引静态路由。那如果我们想支持类似于/hello/:name这样的动态路由怎么办呢？
动态路由，即一条路由规则可以匹配某一类型而非某一条固定的路由。例如/hello/:name，可以匹配/hello/geektutu、hello/jack等。

动态路由有很多种实现方式，支持的规则、性能等有很大的差异。例如开源的路由实现gorouter支持在路由规则中嵌入正则表达式，例如/p/[0-9A-Za-z]+，即路径中的参数仅匹配数字和字母；另一个开源实现httprouter就不支持正则表达式。著名的Web开源框架gin 在早期的版本，并没有实现自己的路由，而是直接使用了httprouter，后来不知道什么原因，放弃了httprouter，自己实现了一个版本。

![](https://geektutu.com/post/gee-day3/trie_eg.jpg)

实现动态路由最常用的数据结构，被称为前缀树(Trie树)。看到名字你大概也能知道前缀树长啥样了：每一个节点的所有的子节点都拥有相同的前缀。这种结构非常适用于路由匹配，比如我们定义了如下路由规则：

/:lang/doc
/:lang/tutorial
/:lang/intro
/about
/p/blog
/p/related
我们用前缀树来表示，是这样的。

![](https://geektutu.com/post/gee-day3/trie_router.jpg)

HTTP请求的路径恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。我们通过树结构查询，如果中间某一层的节点都不满足条件，那么就说明没有匹配到的路由，查询结束。

接下来我们实现的动态路由具备以下两个功能。
1. 参数匹配:。例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
2. 通配*。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器，能够递归地匹配子路径。

# the fourth day 11.7

## TASK
实现路由分组控制(Route Group Control)，代码约50行

## 分组的意义
在真实的业务场景中，往往某一路由需要相似的处理，例如 `/admin`开头的路由往往需要明确访问者身份，而大部分情况下的路由分组是根据prefix来分组的
这里实现的分组控制是以前缀来区分的并且支持分组的嵌套，作用在分组上的中间件（middleware）也会作用在子分组上，子分组也可有自己的中间件，有点像继承
所以我们定义的分组里面必须包含这几项：prefix、parent、middleware，
但还需要访问`router`东西，也就是保存一个指针，指向Engine，整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了。

# the fifth day 11.9

## TASK
1. 设计并实现 Web 框架的中间件(Middlewares)机制。
2. 实现通用的Logger中间件，能够记录请求到响应所花费的时间，代码约50行

## 什么是中间件
中间件（Middleware）是一类软件，它位于客户端和服务器端之间，用来协调和管理客户端和服务器端之间的交互。中间件接收来自客户端的请求，处理请求，然后将结果返回给客户端。中间件可以提供一些额外的功能，例如身份验证、日志记录、缓存、负载均衡等。常见的中间件包括Web服务器、应用服务器、消息队列、数据库连接池等。
非业务的技术类组件，插入点不能太近，会导致功能不足，插入点太底层则会导致逻辑复杂
中间件的输入是什么？中间件的输入，决定了扩展能力。暴露的参数太少，用户发挥空间有限。
由于作者理解的不是很清楚，这里贴出此链接，以供扩展：[一文读懂中间件](https://cloud.tencent.com/developer/article/2029175)

## 关于Next() 
这个函数是一个迭代函数，它在上下文中将处理程序传递给下一个处理程序，直到达到处理程序列表的末尾。具体步骤如下： 
1. 将当前上下文的索引加1，表示移动到下一个处理程序。 
2. 获取处理程序列表的长度。 
3. 通过循环将处理程序传递给当前上下文的索引处的处理程序，直到达到处理程序列表的长度。 
循环中的具体步骤如下： 
1. 将当前处理程序(c.handlers[c.index])传递给当前上下文(c)。 
2. 将当前上下文的索引加1，表示移动到下一个处理程序。 循环会一直执行，直到当前上下文的索引达到处理程序列表的长度。



# the sixth day 11.10

## TASK
1. 实现静态资源服务(Static Resource)。
2. 支持HTML模板渲染。

## html/template 库
html/template为 HTML 提供了较为完整的支持。包括普通变量渲染、列表渲染、对象渲染等。gee 框架的模板渲染直接使用了html/template提供的能力。
```go
type Engine struct {
	...
	//将模板加载进内存
	//是一个指向预定义的HTML模板的指针。可以使用该指针来操作和渲染HTML模板，以便在Go程序中生成动态内容。
	//它是一个template.Template类型的指针，该类型包含了已经定义好的模板和相关的数据，可以用来快速生成HTML页面。
	htmlTemplates *template.Template
	//自定义模板渲染函数
	//"funcMap" 是一个 Go 语言中的 template 包下的函数映射，它是一个 map 类型的变量，用于存储自定义的模板函数以及其对应的处理函数。
	//通过 funcMap，可以在模板中调用自定义的函数并对模板进行动态处理。
	funcMap       template.FuncMap
}
```

# the seventh day 11.11

## TASK
实现错误处理机制

## 为什么要实现错误处理机制
1. 防止服务器宕机
2. 方便用户定位错误

## defer
panic 会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer 的任务，执行完成后再退出。效果类似于 java 语言的 try...catch。

## recover
Go 语言还提供了 recover 函数，可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效。
