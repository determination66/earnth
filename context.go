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

	RespStatusCode int
	ResData        []byte

	pathParams map[string]string

	queryValues url.Values
}

func newContext(req *http.Request, resp http.ResponseWriter) *Context {
	return &Context{
		Req:  req,
		Resp: resp,
	}
}

func (ctx *Context) JSON(statusCode int, data interface{}) error {
	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.RespStatusCode = statusCode
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.ResData = jsonData
	return err
}

func (ctx *Context) pathValue(key string) (string, error) {
	res, ok := ctx.pathParams[key]
	if !ok {
		return "", errors.New("no such path param")
	}
	return res, nil
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

// QueryValue return query value
// add cache
func (ctx *Context) QueryValue(key string) (string, error) {
	if ctx.queryValues == nil {
		ctx.queryValues = ctx.Req.URL.Query()
	}
	vals, ok := ctx.queryValues[key]
	if !ok {
		return "", errors.New("no such query param")
	}
	return vals[0], nil
}
