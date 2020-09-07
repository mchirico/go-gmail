package messages

import (
	"bytes"
	"fmt"
	"github.com/mchirico/go-pubsub/pubsub"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"strings"
	"testing"
	"time"
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
	for id := range r {
		fmt.Println(r[id]["Subject"])
		fmt.Println(r[id]["Message-ID"])
		fmt.Println(r[id]["Return-Path"])
		fmt.Println(r[id]["From"])
		fmt.Println(r[id]["Snippet"])
		fmt.Println(r[id]["Id"])
		fmt.Println("----------------------")

	}

}

func TestGetRaw(t *testing.T) {
	r := GetRaw("TRASH", 1)
	for k, v := range r {
		fmt.Println(k, string(v))
	}
}

func Test_ReturnDomains(t *testing.T) {

	r := GetNewMessages("SPAM", 100)
	Domains(r)
	//	ioutil.WriteFile("domainsBlock",[]byte(st),0644)

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
Please note:  I'm only open to a corp-to-corp contract, 
with my company CWXSTAT INC. Remote contract work ONLY.  

Hourly rate range $98/hr to $117/hr.

Please confirm this position is 100% remote, including 
after COVID-19, and will work on a corp-to-corp contract, 
within the hourly range stated above.  

Please be sure to include a mobile phone number, where you
can be reached by text. I'm sharing my mobile number below.   

If you agree with these conditions (hourly range and remote work), 
we can explore the next step. 

Please confirm and respond to the email ONLY if you 
AGREE with ALL of these conditions.


Regards,

Mike Chirico
mc@cwxstat.com
(215) 326-9389 (text only)`

	subject := "C2C Contracts Only...  Re: " + r[id]["Subject"]
	msgID := r[id]["Message-ID"]
	m, err := Reply(r[id]["Id"], msgID, "mc@cwxstat.com",
		r[id]["From"], subject, msg)
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
	t.Log(m)
}

func TestThread(t *testing.T) {
	Thread("TRASH", 3)
}

func OnlyDoOnce() {

	b, err := ioutil.ReadFile("../../credentials/topic_name.json")

	topic := strings.TrimSuffix(string(b), "\n")
	watchReq := &gmail.WatchRequest{
		LabelIds:  []string{"TRASH"},
		TopicName: topic,
	}

	c := Watch("me", watchReq)
	wr, err := c.Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(wr.Expiration)

	// Convert the milli seconds into seconds
	secs := wr.Expiration / 1000

	tt := time.Unix(secs, 0)
	nanos := wr.Expiration * 1000000
	tM := time.Unix(0, nanos)
	fmt.Printf("Expiration: %s\n", tt)
	fmt.Printf("Expiration2: %s\n", tM)
	fmt.Printf("HistoryId: %d\n", wr.HistoryId)

}

func TestWatch(t *testing.T) {

	g := pubsub.NewG()
	var buf bytes.Buffer

	for i := 0; i < 10; i++ {
		msg, err := g.PullMsgs(&buf, "gmail-sub")
		if err != nil {
			t.Fatalf("No message")
		}
		fmt.Printf("msg: %s\n", msg)
	}

}

func Test_StopWatch(t *testing.T) {
	err := StopWatch("me")
	t.Log(err)
}

func Test_jsonR(t *testing.T) {
	b, err := ioutil.ReadFile("../../credentials/topic_name.json")
	if err != nil {
		fmt.Println(string(b))
	}
	fmt.Printf("->%s<-\n", (string(b)))
}

func Test_play(t *testing.T) {
	OnlyDoOnce()
}

func TestStartWatch(t *testing.T) {
	b, _ := ioutil.ReadFile("../../credentials/topic_name.json")

	topic := strings.TrimSuffix(string(b), "\n")

	StartWatch("me", topic)
}
