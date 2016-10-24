package main

import (
	"fmt"
	"net/http"
	"./lib"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {
	lib.Serve(hello)
}
