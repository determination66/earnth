package earnth

import (
	"fmt"
	"testing"
)

func TestServer(T *testing.T) {
	s := NewHTTPServer()
	//fmt.Println(s)

	s.Get("/user", func() {
		fmt.Println("hello world")
	})

	s.Start(":8080")
}
