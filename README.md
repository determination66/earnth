## 使用go语言编写一个轻量TCP的web框架
> - 仓库地址: https://gitee.com/determination66/earnth
> - 简易go web框架的调研: https://www.kdocs.cn/l/cm4wvF5iq20x
> - go web小白之路: https://www.kdocs.cn/l/cqmUIDERGRa1

## earnth是一个功能全面，开发效率高的http框架

### 快速开始
```
package main
import (
	"fmt"
	"github.com/determination66/earnth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// MyMiddleware 自定义中间件，例如计算耗时
func MyMiddleware(next earnth.HandleFunc) earnth.HandleFunc {
	return func(ctx *earnth.Context) {
		start := time.Now()

		next(ctx)

		// 计算并输出中间件处理时间（微秒）
		duration := time.Since(start)
		log.Printf("\nYour Request took: %v ns\n", duration.Nanoseconds())
	}
}

// 上传

func main() {
	s := earnth.Default()

	s.Get("/users/me", func(ctx *earnth.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"hello": "world",
			"msg":   "ok",
		})
	})

	s.Get("/users/me/:id", func(ctx *earnth.Context) {
		id, _ := ctx.PathValue("id")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"your_id": id,
		})
	})

	// conflict route
	//s.Get("/users/me/*/detail", func(ctx *earnth.Context) {
	//	ctx.JSON(http.StatusOK, map[string]interface{}{
	//		"/users/me/*/detail": "called",
	//	})
	//})

	// 模拟注册，需要客户端传入username和password
	// 传入账号dcl 密码123

	s.Post("/user/login", func(ctx *earnth.Context) {
		time.Sleep(10 * time.Millisecond)
		user, err := ctx.ParseForm("username")
		if err != nil {
			panic(err)
		}
		pwd, err := ctx.ParseForm("password")
		if err != nil {
			panic(err)
		}
		fmt.Println(user, pwd)
		// can fix JSON's RespData
		if user == "dcl" && pwd == "123" {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"msg": "Register succeed!",
			})
		} else {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"msg": "Register failed!",
			})
		}

	})

	// 下载操作
	// http://127.0.0.1:9999/download?file=Jay-sunny%20weather.mp3
	fu := earnth.NewFileDownload(filepath.Join("static", "download"))
	s.Get("/download", fu.Handle())

	//http://127.0.0.1:6666/metrics
	s.Use((&earnth.MiddlewareBuilder{
		Subsystem: "earnth",
		Name:      "http_request",
		Help:      "Help",
		ConstLabels: map[string]string{
			"instance_id": "1234567",
		},
	}).Build(), MyMiddleware)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hello")
		})
		// 一般来说，在实际中我们都会单独准备一个端口给这种监控
		err := http.ListenAndServe(":6666", nil)
		if err != nil {
			fmt.Println("Failed to start server:", err)
		}
	}()

	s.Start(":9999")

}
```