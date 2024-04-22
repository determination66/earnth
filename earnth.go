package earnth

import (
	"fmt"
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	Use(mdls ...MiddlewareFunc)
	// AddRoute 也就是说可以用GET、POST、DELETE、OPTIONS、PUT、TRACE、CONNECT、HEAD
	AddRoute(method string, path string, handleFunc HandleFunc)
}

// HTTPServer This is the earnth's Engine. It exposes all the interfaces for users.
type HTTPServer struct {
	*router

	mdls []MiddlewareFunc

	tplEngine TemplateEngine
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (H *HTTPServer) registerTemplateEngine(tplEngine TemplateEngine) {
	H.tplEngine = tplEngine
}

func (H *HTTPServer) Use(mdls ...MiddlewareFunc) {
	H.mdls = append(H.mdls, mdls...)
}

// ServeHTTP This implements the handler interface,so the earnth's real processing logic code.
func (H *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(req, w)

	var root HandleFunc = H.serve
	// exec mdls
	// callback func need to reverse order
	for i := len(H.mdls) - 1; i >= 0; i-- {
		root = H.mdls[i](root)
	}
	// use m to write back
	var m MiddlewareFunc = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			//fmt.Println("Before m")
			next(ctx)
			//fmt.Println("After m")
			H.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (H *HTTPServer) serve(ctx *Context) {
	dst := H.matchRouter(ctx.Req.Method, ctx.Req.URL.Path)
	// do not match HandleFunc
	if dst == nil || dst.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("404 page not found"))
		return
	}
	// put the matchInfo into ctx
	ctx.pathParams = dst.pathParams
	ctx.MatchedRoute = dst.n.routerPath
	dst.n.handler(ctx)
}

func (H *HTTPServer) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Println("Listening on " + addr)
	return http.Serve(ln, H)
}

func (H *HTTPServer) Get(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodGet, path, handleFunc)
}

func (H *HTTPServer) Post(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodPost, path, handleFunc)
}

func (H *HTTPServer) Options(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodOptions, path, handleFunc)
}

func (H *HTTPServer) Delete(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodDelete, path, handleFunc)
}

func (H *HTTPServer) Put(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodPut, path, handleFunc)
}

func (H *HTTPServer) Patch(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodPatch, path, handleFunc)
}

func (H *HTTPServer) Head(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodHead, path, handleFunc)
}
func (H *HTTPServer) Trace(path string, handleFunc HandleFunc) {
	H.AddRoute(http.MethodTrace, path, handleFunc)
}

func (H *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 && !ctx.RespHeaderCommitted {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		panic("fail to write back")
	}
}
