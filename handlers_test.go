package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("NewRequest: %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(rootHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d.\n Got %d instead.", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), "Enter trading symbol here...") {
		t.Errorf("Unexpected response body: %s", rr.Body.String())
	}

}

func TestSearchHandler(t *testing.T) {
	tres := `{ success:true ,results:1,rows:[{Symbol:"S",CompanyName:"CN",ISIN:"ISIN123",Ind:"-",Purpose:"Purpose",BoardMeetingDate:"20-May-2020",DisplayDate:"24-Mar-2020",seqId:"12345678",Details:"Details"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tres)
	}))
	defer ts.Close()

	xchg = &exchange{URL: ts.URL}

	body := []byte("searchTerm=FOO")
	req, err := http.NewRequest("POST", "/search", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	http.HandlerFunc(searchHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d.\n Got %d instead.", http.StatusOK, status)
	}

	tmpl, err := template.ParseFiles("result.html")
	if err != nil {
		t.Errorf("template.ParseFiles: %v", err)
	}

	m := &Meeting{"20-May-2020", "Purpose", "Details"}

	var expected bytes.Buffer
	if err := tmpl.Execute(&expected, m); err != nil {
		t.Errorf("tmpl.Execute: %v", err)
	}

	if got := rr.Body.String(); got != expected.String() {
		t.Errorf("Unexpected response body: %s", got)
	}

}
