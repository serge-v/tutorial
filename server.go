// tutorial project contains tutorials for my wife to start with Go language.
// The goal is to teach her to
//  - create simple web apps
//  - deploy apps to the cloud web server
//  - generate html using templates
//  - process html forms
//  - interact with mariadb
//  - create browser notifications in javascript
//
// blog program itself is a CGI handler for golang.org/x/tools/blog package.
// Locally it runs under debug-server.
// In production it runs under Apache web server.
package main

import (
	"net/http/cgi"
//	"golang.org/x/tools/blog"
	"./blogm" // copy of golang.org/x/tools/blog. Used for debugging with printls.
)

var conf = blog.Config {
	ContentPath: "content",
	TemplatePath: "template",
	BaseURL: "http://localhost",
	BasePath: "/blog",
	GodocURL: "http://golang.org",
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
