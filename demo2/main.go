package main

import (
	"earnth/demo2/gee"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

//func onlyForV2() gee.HandlerFunc {
//	return func(c *gee.Context) {
//		// Start timer
//		t := time.Now()
//		// if a server error occurred
//		//c.Fail(500, "Internal Server Error")
//		// Calculate resolution time
//		log.Printf("----------[%d] %s in %v for group v2\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
//	}
//}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	// 文件映射
	r.Static("/assets", "./demo2/static")
	r.GET("/index", func(c *gee.Context) {
		c.HTMLString(http.StatusOK, "<h1>Index Page</h1>")
	})

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("./demo2/templates/*")
	r.GET("/", func(c *gee.Context) {
		c.HTMLTemplate(http.StatusOK, "css.tmpl", nil)
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTMLString(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := v1.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	v1.GET("/hhh", func(c *gee.Context) {
		c.HTMLString(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.Run(":9999")
}
