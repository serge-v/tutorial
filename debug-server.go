// debug-server runs blog locally on http://localhost:8081.
// Start it and edit articles.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func main() {

	fmt.Println("starting blog on http://localhost:8081")

	home := cgi.Handler{}
	home.Path = "./wethome"

	blog := cgi.Handler{}
	blog.Path = "./blog"

	mux := http.NewServeMux()
	mux.Handle("/", &home)
	mux.Handle("/blog/", &blog)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}
