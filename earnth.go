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

	// AddRoute 也就是说可以用GET、POST、DELETE、OPTIONS、PUT、TRACE、CONNECT、HEAD
	AddRoute(method string, path string, handleFunc HandleFunc)
}

// HTTPServer This is the earnth's Engine. It exposes all the interfaces for users.
type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP This implements the handler interface,so the earnth's real processing logic code.
func (H *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 处理业务。
	//ctx := &Context{
	//	Req:    req,
	//	Writer: w,
	//}
	H.ctx = &Context{
		Req:  req,
		Resp: w,
	}

	H.serve()
}

func (H *HTTPServer) serve() {
	dst := H.matchRouter(H.ctx.Req.Method, H.ctx.Req.URL.Path)
	// do not match HandleFunc
	if dst == nil || dst.n.handler == nil {
		H.ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = H.ctx.Resp.Write([]byte("404 page not found"))
		return
	}
	// put the matchInfo into ctx
	H.ctx.pathParams = dst.pathParams
	dst.n.handler(H.ctx)
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
