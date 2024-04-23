package earnth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RecoveryMiddlewareBuilder struct {
	recoveryFunc func(ctx *Context, err any)
}

func Recovery() MiddlewareFunc {
	return NewRecoveryMiddlewareBuilder().Build()
}

func NewRecoveryMiddlewareBuilder() *RecoveryMiddlewareBuilder {
	return &RecoveryMiddlewareBuilder{
		recoveryFunc: func(ctx *Context, err any) {
			// Serialize error information into JSON format
			recoveryLog := struct {
				Host       string `json:"host"`
				Path       string `json:"path"`
				HTTPMethod string `json:"http_method"`
				Error      string `json:"error"`
			}{
				Host:       ctx.Req.Host,
				Path:       ctx.Req.URL.Path,
				HTTPMethod: ctx.Req.Method,
				Error:      fmt.Sprintf("%v", err),
			}
			val, _ := json.Marshal(recoveryLog)

			log.Printf("Panic occurred: %v", err)
			log.Printf("Recovery Log: %s\n", val)

			// You can customize the recovery behavior here, such as sending error response to the client
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.Resp.WriteHeader(ctx.RespStatusCode)
			ctx.RespHeaderCommitted = true
			ctx.Resp.Write([]byte("Internal Server Error,statusCode:500"))
		},
	}
}

func (r *RecoveryMiddlewareBuilder) RegisterLogFunc(f func(ctx *Context, err any)) *RecoveryMiddlewareBuilder {
	r.recoveryFunc = f
	return r
}

func (r *RecoveryMiddlewareBuilder) Build() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					r.recoveryFunc(ctx, err)
				}
				//}else {
				//		ctx.RespStatusCode = http.StatusOK
				//	}
			}()
			//fmt.Println("Before recovery")
			next(ctx)
			//fmt.Println("After recovery")
		}
	}
}
