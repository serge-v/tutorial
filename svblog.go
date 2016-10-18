package main

import (
//	"net/http"
	"net/http/cgi"
//	"golang.org/x/tools/blog"
	"./blogm"
//	"fmt"
)

var conf = blog.Config {
	ContentPath: "content",
	TemplatePath: "template",
	BaseURL: "http://localhost",
	BasePath: "/blog",
	GodocURL: "",
	Hostname: "hostname",
	HomeArticles: 5,
	FeedArticles: 5,
	FeedTitle: "feed",
	PlayEnabled: true,
}

func main() {
	server, err := blog.NewServer(conf)
	if err != nil {
		panic(err)
	}

	err = cgi.Serve(server)
	if err != nil {
		panic(err)
	}
}
