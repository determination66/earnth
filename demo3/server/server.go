package server

import (
	"fmt"
	"net/http"
)

type Server struct {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ServerHttp")
}
