package main

import (
	"html/template"
	"net/http"
	"regexp"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PostFormValue("searchTerm")

	if !regexp.MustCompile(`\A[A-Z0-9]+\z`).MatchString(searchTerm) {
		http.Error(w, "Bad symbol: "+searchTerm, http.StatusBadRequest)
		return
	}

	meetingData, err := xchg.fetchMeeting(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	meeting := parseMeeting(meetingData)

	tmpl, err := template.ParseFiles("result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, meeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
