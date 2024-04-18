package server

import (
	"net/http"
	"testing"
)

func TestServer(T *testing.T) {
	http.ListenAndServe(":8081", &Server{})
}
