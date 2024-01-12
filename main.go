package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello,world!")
	})

	fmt.Println("serve is listening...")
	http.ListenAndServe(":8080", nil)

}
