package main

import (
	"fmt"
	"net/http"
	"log"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {
	fmt.Println("=== starting server on http://localhost:8085/ ===")
	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
