package main

import (
	"fmt"
	"net/http"
	"time"
	"./lib"
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

func getTimes(w http.ResponseWriter, req *http.Request) {
	printTime(w, "GMT")
	printTime(w, "America/New_York")
	printTime(w, "America/Los_Angeles")
	printTime(w, "Europe/Minsk")
	printTime(w, "Europe/Tallinn")
	printTime(w, "Europe/Moscow")
}

func main() {
	lib.Serve(getTimes)
}
