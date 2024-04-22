package recovery

import (
	"earnth"
	"earnth/middlewares/logger"
	"net/http"
	"testing"
)

func TestRecoveryMiddlewareBuilder_Build(t *testing.T) {
	s := earnth.NewHTTPServer()

	s.Use(logger.Logger(), Recovery())

	s.Get("/order/detail", func(ctx *earnth.Context) {
		panic("报错")
		ctx.Resp.Write([]byte("order detail"))
		ctx.Resp.Write([]byte("order detail 2"))
	})
	s.Get("/user/:id", func(ctx *earnth.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"msg": "hello",
		})
	})

	s.Start(":9999")
}
