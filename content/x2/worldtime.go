package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"log"
	"time"
	"os"
)

func printTime(w http.ResponseWriter, timezone string) {
	now := time.Now() // get current time

	zone, err := time.LoadLocation(timezone) // load time zone info
	if err != nil {                          // check for error
		fmt.Println(err.Error())         // print the error if invalid timezone
	} else {
		fmt.Fprintf(w, "%-20s:  %s\n", timezone, now.In(zone))
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	printTime(w, "GMT")
	printTime(w, "America/New_York")
	printTime(w, "America/Los_Angeles")
	printTime(w, "Europe/Minsk")
	printTime(w, "Europe/Tallinn")
	printTime(w, "Europe/Moscow")
}

func main() {
	protocol := os.Getenv("SERVER_PROTOCOL")

	if len(protocol) == 0 {
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
