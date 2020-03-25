package main

import "testing"

func TestParseMeeting(t *testing.T) {
	tests := []struct {
		data    string
		meeting *Meeting
	}{
		{
			``,
			&Meeting{"-", "-", "-"}},
		{
			`{ success:true ,results:1,rows:[{Symbol:"S",CompanyName:"CN",ISIN:"ISIN123",Ind:"-",Purpose:"Purpose",BoardMeetingDate:"20-May-2020",DisplayDate:"24-Mar-2020",seqId:"12345678",Details:"Details"}]}`,
			&Meeting{"20-May-2020", "Purpose", "Details"},
		},
	}

	for _, test := range tests {
		m := parseMeeting(test.data)
		if *m != *test.meeting {
			t.Errorf("Expected %v. Got %v instead", test.meeting, m)
		}
	}

}
