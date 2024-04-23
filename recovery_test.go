package earnth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRecoveryMiddlewareBuilder_Build(t *testing.T) {
	s := NewHTTPServer()

	s.Use(Logger(), Recovery())

	s.Get("/order/detail", func(ctx *Context) {
		panic("报错")
		fmt.Println("hello world")
		ctx.Resp.Write([]byte("order detail"))
		ctx.Resp.Write([]byte("order detail 2"))
	})
	s.Get("/user/:id", func(ctx *Context) {
		fmt.Println("123")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"msg": "hello",
		})
	})

	s.Start(":9999")
}
