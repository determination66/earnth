package earnth

import (
	"net/http"
	"testing"
)

func TestLoggerMiddlewareBuilder_Build(t *testing.T) {
	s := NewHTTPServer()

	s.Use(Logger())

	s.Post("/order/detail", func(ctx *Context) {
		panic("报错")
		ctx.Resp.Write([]byte("order detail"))
	})
	s.Get("/user/:id", func(ctx *Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"msg": "hello",
		})
	})

	s.Start(":9999")
}
