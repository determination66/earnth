package logger

import (
	"earnth"
	"net/http"
	"testing"
)

func TestLoggerMiddlewareBuilder_Build(t *testing.T) {
	s := earnth.NewHTTPServer()

	s.Use(Logger())

	s.Post("/order/detail", func(ctx *earnth.Context) {
		panic("报错")
		ctx.Resp.Write([]byte("order detail"))
	})
	s.Get("/user/:id", func(ctx *earnth.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"msg": "hello",
		})
	})

	s.Start(":9999")
}