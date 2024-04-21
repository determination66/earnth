package logger

import (
	"earnth"
	"encoding/json"
	"log"
)

// LoggerMiddlewareBuilder  for log
type LoggerMiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func Logger() earnth.MiddlewareFunc {
	return NewLoggerMiddlewareBuilder().Build()
}

func (b *LoggerMiddlewareBuilder) SetLogFunc(logFunc func(accessLog string)) *LoggerMiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

func NewLoggerMiddlewareBuilder() *LoggerMiddlewareBuilder {
	return &LoggerMiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

// AccessLog 包含访问日志的信息
type AccessLog struct {
	Host string `json:"host"`
	//Route      string `json:"route"`
	StatusCode int    `json:"status_code"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path"`
}

// Build 用于构建日志记录中间件
func (b *LoggerMiddlewareBuilder) Build() earnth.MiddlewareFunc {
	return func(next earnth.HandleFunc) earnth.HandleFunc {
		return func(ctx *earnth.Context) {
			defer func() {
				l := AccessLog{
					Host: ctx.Req.Host,
					//Route:      "todo",
					StatusCode: ctx.RespStatusCode,
					Path:       ctx.Req.URL.Path,
					HTTPMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				// the accessLog text
				b.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}
