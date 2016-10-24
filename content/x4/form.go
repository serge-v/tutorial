package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {

	comments := loadCommentsFile()
	article := loadArticleFile()

	if r.Method == "GET" {
		article = insertComments(article, comments)
		fmt.Fprintf(w, article)
		return
	}

	// process new comment from form

	comment, author := parseRequest(r)
	if len(comment) > 0 && len(author) > 0 {
		comments = appendComment(comments, comment, author)
		saveCommentsFile(comments)
	}

	http.Redirect(w, r, r.URL.Path, http.StatusFound)
}

func main() {

	protocol := os.Getenv("SERVER_PROTOCOL")

	if len(protocol) == 0 {
		fmt.Println("=== starting server on http://localhost:8085/ ===")
		http.HandleFunc("/", requestHandler)
		log.Fatal(http.ListenAndServe(":8085", nil))
		return
	}

	err := cgi.Serve(http.HandlerFunc(requestHandler))
	if err != nil {
		panic(err)
	}
}

func parseRequest(r *http.Request) (string, string) {
	if err := r.ParseForm(); err != nil {
		return "error", "error"
	}

	var author, comment string

	if len(r.Form["author"]) > 0 {
		author = r.Form["author"][0]
	}

	if len(r.Form["comment"]) > 0 {
		comment = r.Form["comment"][0]
	}
	
	if len(author) > 20 {
		author = author[:20]
	}

	if len(comment) > 150 {
		comment = comment[:150]
	}

	return comment, author
}

func loadCommentsFile() string {
	bytes, err := ioutil.ReadFile("comments.txt")
	if err != nil {
		return "" // no file, just return empty string
	}
	return string(bytes)
}

func appendComment(text, comment, author string) string {
	text += html.EscapeString(comment) + " (<i>by " + html.EscapeString(author) + "</i>)<br><br>\n"
	return text
}

func loadArticleFile() string {
	bytes, err := ioutil.ReadFile("article.html")
	if err != nil {
		return "" // no file, just return empty string
	}
	return string(bytes)
}

func saveCommentsFile(text string) {
	err := ioutil.WriteFile("comments.txt", []byte(text), 0666)
	if err != nil {
		panic(err)
	}
}

func insertComments(article, comments string) string {
	s := strings.Replace(article, "{placeholder_for_comments}", comments, 1)
	return s
}
