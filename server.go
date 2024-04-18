package earnth

import "net/http"

type HandleFunc func()

type Server interface {
	http.Handler
	Start()

	// AddRoute 也就是说可以用GET、POST、DELETE、OPTIONS、PUT、TRACE、CONNECT、HEAD
	AddRoute(method string, path string, handleFunc HandleFunc)
}

type HTTPServer struct {
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (H *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (H *HTTPServer) Start() {
	//TODO implement me
	panic("implement me")
}

func (H *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
	//TODO implement me
	panic("implement me")
}
