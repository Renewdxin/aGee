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
