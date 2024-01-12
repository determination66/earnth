package main

import (
	"net/http"
)

func main() {
	//http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
	//	fmt.Fprintf(w, "hello,world!")
	//})
	//
	//fmt.Println("serve is listening...")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}
