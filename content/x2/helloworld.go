package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"log"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {
	if false {
		fmt.Println("=== starting server on http://localhost:8085/ ===")
		http.HandleFunc("/", hello)
		log.Fatal(http.ListenAndServe(":8085", nil))
	else {
		err := cgi.Serve(http.HandlerFunc(hello))
		if err != nil {
			panic(err)
		}
	}
}
