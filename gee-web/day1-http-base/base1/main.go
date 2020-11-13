package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	//HandleFunc注册给定模式的处理函数
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	//ListenAndServer侦听tcp网络地址addr.  然后调用预处理程序
	//一起使用. 以处理传入连接上的请求. 将接受的连接配置为启用 TCP保持活动连接
	//处理程序通常为nil.在这种情况下. 将使用DefaultServerMux
	//listenAndServe总是返回非nil错误
	log.Fatal(http.ListenAndServe(":9999", nil))

}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	//Fprintf根据格式说明符格式化并写入W.
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
