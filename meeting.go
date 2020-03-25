package main

import "regexp"

type Meeting struct {
	Date    string
	Purpose string
	Details string
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
