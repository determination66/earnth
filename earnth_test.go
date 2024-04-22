package earnth

import (
	"fmt"
	"net/http"
	"testing"
	"time"
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
		name, _ := ctx.pathValue("name")
		ctx.Resp.Write([]byte(name))
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

	http.HandleFunc("/query1", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 多次调用 Query() 方法
		for i := 0; i < 5; i++ {
			queryParams := r.URL.Query()
			username := queryParams.Get("username")
			fmt.Fprintf(w, "Username: %s\n", username)
		}

		elapsed := time.Since(start)
		fmt.Fprintf(w, "Time elapsed: %s\n", elapsed)
	})

	http.ListenAndServe(":9999", nil)
}

func TestUse(T *testing.T) {

	s := NewHTTPServer()
	mdls := []MiddlewareFunc{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个Before")
				next(ctx)
				fmt.Println("第一个After")

			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个Before")
				next(ctx)
				fmt.Println("第二个After")
			}
		},
	}
	//tpls, err := template.ParseGlob("testdata/tpls/*.gohtml")
	//if err != nil {
	//	panic(err)
	//}
	//s.registerTemplateEngine(&GoTemplateEngine{
	//	T: tpls,
	//})

	err := s.LoadGlob("testdata/tpls/*.gohtml")
	if err != nil {
		panic(err)
	}
	s.Use(mdls...)

	s.Get("/login", func(ctx *Context) {
		err = ctx.Render("login.gohtml", map[string]interface{}{"email": "2496417370@qq.com", "password": "123456"})
		if err != nil {
			panic(err)
		}
	})
	s.Post("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("order detail"))
	})

	s.Start(":9999")
}
