package earnth

import (
	"fmt"
	"testing"
)

func TestServer(T *testing.T) {
	s := NewHTTPServer()
	//fmt.Println(s)

	s.Get("/user", func(ctx *Context) {
		fmt.Println("hello world")
	})
	s.Post("/order/detail", func(ctx *Context) {
		ctx.Writer.Write([]byte("order detail"))
	})

	s.Start(":8080")
}
