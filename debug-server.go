// debug-server OMIT
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func main() {

	fmt.Println("starting blog on http://localhost:8081/blog")

	blog := cgi.Handler{}
	blog.Path = "./blog"

	mux := http.NewServeMux()
	mux.Handle("/blog/", &blog)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}
// debug-server OMIT
