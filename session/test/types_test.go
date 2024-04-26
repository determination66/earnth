package test

import (
	"github.com/determination66/earnth"
	earnth2 "github.com/determination66/earnth/session"
	"net/http"
	"testing"
)

// 登录校验
var p earnth2.Propagator
var s earnth2.Store

func LoginMiddleware(next earnth.HandleFunc) earnth.HandleFunc {
	return func(ctx *earnth.Context) {
		if ctx.Req.URL.Path == "/login" {
			// 放过，用户登录
			next(ctx)
			return
		}
		sessId, err := p.Extract(ctx.Req)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("please login again")
			return
		}
		_, err = s.Get(ctx.Req.Context(), sessId)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("please login again")
			return
		}
	}
}

// 原始session的demo
func TestSession(t *testing.T) {
	server := earnth.NewHTTPServer()

	server.Use(LoginMiddleware)
	server.Get("/user", func(ctx *earnth.Context) {
		sessId, err := p.Extract(ctx.Req)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("please login again")
			return
		}

		sess, err := s.Get(ctx.Req.Context(), sessId)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("please login again")
			return
		}
		nike, _ := sess.Get(ctx.Req.Context(), "nikename")

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"msg":      "ok",
			"nikename": nike,
		})
	})
	server.Start(":9999")
}
