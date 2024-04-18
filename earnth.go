package earnth

import (
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
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (H *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 处理业务。

	//fmt.Println("serveHTTP")
}

func (H *HTTPServer) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
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
