package logger

import (
	"earnth"
	"testing"
)

func TestLoggerMiddlewareBuilder_Build(t *testing.T) {
	s := earnth.NewHTTPServer()

	s.Use(Logger())

	s.Post("/order/detail", func(ctx *earnth.Context) {
		ctx.Resp.Write([]byte("order detail"))
	})

	s.Start(":9999")
}
