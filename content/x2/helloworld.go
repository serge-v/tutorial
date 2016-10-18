package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"log"
	"flag"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world. local = %v\n", *local)
}

var local = flag.Bool("local", false, "start local server for debugging")

func main() {
	flag.Parse()

	if *local {
		fmt.Println("=== starting server on http://localhost:8085/ ===")
		http.HandleFunc("/", hello)
		log.Fatal(http.ListenAndServe(":8085", nil))
		return
	}
	
	err := cgi.Serve(http.HandlerFunc(hello))
	if err != nil {
		panic(err)
	}
}
