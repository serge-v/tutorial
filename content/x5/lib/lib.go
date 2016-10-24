package lib // name of the package

// I copied all imports from helloworld.go

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"log"
	"os"
)

// I copied Serve from helloworld.go main and added handler http.HandlerFunc parameter

func Serve(handler http.HandlerFunc) {

	protocol := os.Getenv("SERVER_PROTOCOL")

	if len(protocol) == 0 {
		fmt.Println("=== starting server on http://localhost:8085/ ===")
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8085", nil))
		return
	}

	err := cgi.Serve(handler)
	if err != nil {
		panic(err)
	}
}
