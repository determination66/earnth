package test

import (
	"github.com/determination66/earnth"
	"github.com/determination66/earnth/session"
	"github.com/determination66/earnth/session/cookie"
	"github.com/determination66/earnth/session/memory"
	"net/http"
	"testing"
	"time"
)

// 登录校验
var m *session.Manager = &session.Manager{
	Propagator: cookie.NewPropagator(),
	Store:      memory.NewStore(time.Minute * 15),
	CtxSessKey: "sessKey",
}

func LoginMiddleware(next earnth.HandleFunc) earnth.HandleFunc {
	return func(ctx *earnth.Context) {
		if ctx.Req.URL.Path == "/login" {
			// 放过，用户登录
			next(ctx)
			return
		}
		// 中间件
		_, err := m.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("Unauthorized,please login again")
			return
		}
		//插入刷新session的逻辑
		err = m.RefreshSession(ctx)
		if err != nil {
			ctx.RespData = []byte("refresh session failed")
			return
		}

		next(ctx)
	}
}

// 原始session的demo
func TestSession(t *testing.T) {
	server := earnth.Default()

	server.Use(LoginMiddleware)

	server.Post("/login", func(ctx *earnth.Context) {
		// 校验用户名密码

		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusUnauthorized
			ctx.RespData = []byte("please login again")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "xiaoming")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("please login again")
			return
		}
		ctx.RespStatusCode = http.StatusOK

	})

	server.Post("/logout", func(ctx *earnth.Context) {
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("remove session failed")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("logout successfully")
	})

	server.Get("/user", func(ctx *earnth.Context) {
		sess, _ := m.GetSession(ctx)
		// 假如说我要把昵称从 session 里面拿出来
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})
	server.Start(":9999")
}
