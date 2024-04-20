package earnth

import "net/http"

type Context struct {
	Req        *http.Request
	Writer     http.ResponseWriter
	pathParams map[string]string
	statusCode int
}

func (ctx *Context) getParam(key string) string {
	return ctx.pathParams[key]
}
