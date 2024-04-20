package earnth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	pathParams map[string]string

	queryValues url.Values
	statusCode  int
}

func newContext(req *http.Request, resp http.ResponseWriter) *Context {
	return &Context{
		Req:  req,
		Resp: resp,
	}
}

func (ctx *Context) getParam(key string) string {
	return ctx.pathParams[key]
}

// BindJSON the most popular method
func (ctx *Context) BindJSON(obj interface{}) error {
	if ctx.Req.Body == nil {
		return errors.New("body is nil")
	}
	decoder := json.NewDecoder(ctx.Req.Body)
	return decoder.Decode(obj)
}

// ParseForm parse form the form and query
func (ctx *Context) ParseForm(key string) (string, error) {
	err := ctx.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return ctx.Req.Form.Get(key), nil
}

func (ctx *Context) Query(key string) string {
	return ""
	//return ctx.Req.URL.Query()
}
