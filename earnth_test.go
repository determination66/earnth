package earnth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(T *testing.T) {
	s := NewHTTPServer()
	//fmt.Println(s)

	//s.Get("/user/*", func(ctx *Context) {
	//	fmt.Println("hello world")
	//})
	//s.Get("/user/*", func(ctx *Context) {
	//	fmt.Println("hello world")
	//})

	s.Get("/", func(ctx *Context) {
		fmt.Println("/")
	})
	//s.Get("/user", func(ctx *Context) {
	//	fmt.Println("hello world")
	//})
	s.Post("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("order detail"))
	})
	s.Get("/user/:name", func(ctx *Context) {
		ctx.Resp.Write([]byte(ctx.getParam("name")))
	})

	s.Post("/submit", func(ctx *Context) {
		username, _ := ctx.ParseForm("username")
		password, _ := ctx.ParseForm("password")
		// 处理表单数据
		fmt.Fprintf(ctx.Resp, "Received form data:\n")
		fmt.Fprintf(ctx.Resp, "Username: %s\n", username)
		fmt.Fprintf(ctx.Resp, "Password: %s\n", password)
	})

	s.Start(":9999")
}

func TestOrigin(T *testing.T) {
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		// 解析 URL 中的查询参数
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}

		// 从 URL 查询参数中获取参数值
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		// 处理表单数据
		fmt.Fprintf(w, "Received form data:\n")
		fmt.Fprintf(w, "Username: %s\n", username)
		fmt.Fprintf(w, "Password: %s\n", password)
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		// 解析 URL 中的查询参数
		queryParams := r.URL.Query()

		// 从查询参数中获取参数值
		username := queryParams.Get("username")
		password := queryParams.Get("password")

		// 处理查询参数
		fmt.Fprintf(w, "Received query parameters:\n")
		fmt.Fprintf(w, "Username: %s\n", username)
		fmt.Fprintf(w, "Password: %s\n", password)
	})

	http.ListenAndServe(":9999", nil)
}
