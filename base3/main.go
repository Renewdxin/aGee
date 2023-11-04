package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/ping", func(req *http.Request, w http.ResponseWriter) {
		fmt.Fprintf(w, "Header[%q] = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(req *http.Request, w http.ResponseWriter) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")

}
