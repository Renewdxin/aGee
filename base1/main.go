package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// http.HandleFunc将指定的处理函数绑定到根路径上
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	//第二个参数代表处理所有的http请求，nil代表使用标准库中的实例处理
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
