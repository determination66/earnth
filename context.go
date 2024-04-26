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

	RespStatusCode      int
	RespHeaderCommitted bool
	RespData            []byte

	pathParams map[string]string

	MatchedRoute string

	tplEngine TemplateEngine

	queryValues url.Values
	//RespCommitted bool // Add a field to mark if the response has been committed

	UserValues map[string]any
}

func newContext(req *http.Request, resp http.ResponseWriter, tplEngine TemplateEngine) *Context {
	return &Context{
		Req:       req,
		Resp:      resp,
		tplEngine: tplEngine,
	}
}

func (ctx *Context) Render(tplName string, data any) error {
	// not ok
	// tplName = tplName + ".gohtml"
	// tplName = tplName + c.tplPrefix
	var err error
	if ctx.tplEngine == nil {
		panic("ctx template engine is nil, please call RegisterTemplateEngine first")
	}
	ctx.RespData, err = ctx.tplEngine.Render(ctx.Req.Context(), tplName, data)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		return err
	}
	ctx.RespStatusCode = http.StatusOK
	return nil
}

func (ctx *Context) JSON(statusCode int, data interface{}) error {
	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.RespStatusCode = statusCode
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.RespData = append(ctx.RespData, jsonData...)
	return err
}

func (ctx *Context) PathValue(key string) (string, error) {
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
	//decode val
	// find bug
	val := vals[0]
	val, err := url.QueryUnescape(val)
	if err != nil {
		return "", err
	}
	return val, nil
	//return vals[0], nil
}
