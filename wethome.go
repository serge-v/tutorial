// wethome is a home page application for wet.voilokov.com.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

var html = `
<html>
<head>
<title>wet.voilokov.com</title>
<head>
<body>
<a href="/blog/">Golang tutorial</a>
</body>
</html>
`

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, html)
}

func main() {	
	err := cgi.Serve(http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}
