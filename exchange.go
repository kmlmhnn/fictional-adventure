package main

import (
	"io/ioutil"
	"net/http"
)

type exchange struct {
	URL string
}

func (e *exchange) fetchMeeting(symbol string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", e.URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", "Latest Announced")
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
