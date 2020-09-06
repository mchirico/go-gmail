package messages

import (
	"fmt"
	"testing"
)

func TestLabels(t *testing.T) {
	m, err := Labels()
	if err != nil {
		t.Fatalf("Labels: %v\n", err)
	}
	if _, ok := m["TRASH"]; !ok {
		t.Fatalf("Labels: no TRASH label\n")
	}
}

func TestGetNewMessages(t *testing.T) {
	r := GetNewMessages("TRASH", 3)
	id := 0
	fmt.Println(r[id]["Subject"])
	fmt.Println(r[id]["Message-ID"])
	fmt.Println(r[id]["Return-Path"])
	fmt.Println(r[id]["From"])
	fmt.Println(r[id]["Snippet"])
	fmt.Println(r[id]["Id"])

}



func TestGetRaw(t *testing.T) {
	r := GetRaw("TRASH", 1)
	for k, v := range r {
		fmt.Println(k, string(v))
	}
}


func Test_Reply(t *testing.T) {
	r := GetNewMessages("TRASH", 1)
	id := 0
	fmt.Println(r[id]["Subject"])
	fmt.Println(r[id]["Message-ID"])
	fmt.Println(r[id]["Return-Path"])
	fmt.Println(r[id]["From"])
	fmt.Println(r[id]["Snippet"])
	fmt.Println(r[id]["Id"])

	msg := `
Please note:  I'm only open to a corp-to-corp contract, with my company CWXSTAT INC. Remote contract work ONLY.  Hourly rate between $98/hr to $117/hr.

Please confirm this position is 100% remote,including after COVID-19, and will work on a corp-to-corp contract, within the hourly range stated above.  

Please be sure to include a mobile phone number, where you can be reached by text. I'm sharing my mobile number below.   

If you agree with these conditions (hourly range and remote work), we can explore the next step. 

Please confirm and respond to the email ONLY if you AGREE with ALL of these conditions.


Regards,

Mike Chirico
mc@cwxstat.com
(215) 326-9389 (text only)`

Reply(r[id]["Id"][0],"mc@cwxstat.com",r[id]["From"][0],msg)



}
