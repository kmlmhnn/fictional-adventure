package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Meeting struct {
	Date    string
	Purpose string
	Details string
}

func fetchMeeting(symbol string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www1.nseindia.com/corporates/corpInfo/equities/getBoardMeetings.jsp", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", "Latest Announced")
	q.Add("industry", "")
	q.Add("Symbol", symbol)
	q.Add("Period", "Latest Announced")
	q.Add("Industry", "")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func parseMeeting(data string) *Meeting {
	findAttrib := func(pattern string) string {
		matches := regexp.MustCompile(pattern).FindStringSubmatch(data)
		if len(matches) != 2 {
			return "-"
		} else {
			return matches[1]
		}
	}
	m := &Meeting{}
	m.Date = findAttrib(`BoardMeetingDate:"([^"]*)"`)
	m.Purpose = findAttrib(`Purpose:"([^"]*)"`)
	m.Details = findAttrib(`Details:"([^"]*)"`)
	return m
}

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
		http.Error(w, "Bad symbol", http.StatusBadRequest)
		return
	}

	data, err := fetchMeeting(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	meeting := parseMeeting(data)

	tmpl, err := template.ParseFiles("result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, meeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
