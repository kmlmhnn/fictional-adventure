package main

import (
	"log"
	"net/http"
)

var xchg *exchange

func init() {
	xchg = &exchange{
		URL: "https://www1.nseindia.com/corporates/corpInfo/equities/getBoardMeetings.jsp",
	}
}

func main() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
