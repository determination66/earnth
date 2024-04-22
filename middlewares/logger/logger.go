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

func NewLoggerMiddlewareBuilder() *LoggerMiddlewareBuilder {
	return &LoggerMiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog + "\n")
		},
	}
}

// Logger the Default logger
// if you want to type,you can use RegisterLogFunc for your logger
func Logger() earnth.MiddlewareFunc {
	return NewLoggerMiddlewareBuilder().Build()
}

func (b *LoggerMiddlewareBuilder) RegisterLogFunc(logFunc func(accessLog string)) *LoggerMiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

// AccessLog log info
type AccessLog struct {
	Host string `json:"host"`
	//Route      string `json:"route"`
	StatusCode int    `json:"status_code"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path"`
}

// Build for build logger middleware
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
			//fmt.Println("Before loggerr")
			next(ctx)
			//fmt.Println("After loggerr")
		}
	}
}
