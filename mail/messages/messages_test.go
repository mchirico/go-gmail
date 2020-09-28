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
	r, err := GetNewMessages("TRASH", 3)
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
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

	r, _ := GetNewMessages("SPAM", 1000)
	m := Domains(r)
	s := ""
	for k, v := range m {
		s += fmt.Sprintf("%s,%d\n", k, v)
	}
	ioutil.WriteFile("domainsBlock", []byte(s), 0644)

}

func Test_Reply(t *testing.T) {
	r, _ := GetNewMessages("TRASH", 1)
	id := 0
	fmt.Println(r[id]["Subject"])
	fmt.Println(r[id]["Message-ID"])
	fmt.Println(r[id]["Return-Path"])
	fmt.Println(r[id]["From"])
	fmt.Println(r[id]["Snippet"])
	fmt.Println(r[id]["Id"])

	if _, ok := r[id]["AI-Msg-Field"]; ok {
		return
	}

	msg := `
Please note:  I'm only open to a corp-to-corp contract, 
with my company CWXSTAT INC. Remote contract work ONLY.  

Is this for a contract?


Regards,

Mike Chirico
mc@cwxstat.com
(215) 326-9389 (text only)`

	subject := "Contract? Remote? Re: " + r[id]["Subject"]
	msgID := r[id]["Message-ID"]
	m, err := ReplyAI(r[id]["Id"], msgID, "mc@cwxstat.com",
		r[id]["From"], subject, msg, "contract")
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
	t.Log(m)
}

func TestSendContentType(t *testing.T) {
	msg := `stuff`
	r := SendContentType("mchirico@gmail.com",
		"test1", msg)
	//headers := r.Header()
	//value := `multipart/alternative; boundary="_=_swift-6292908865f5a34286af589.42593834_=_"`
	//headers.Set("Subject","bozo")

	fmt.Println(r)
	message, err := r.Do()
	if err != nil {
		fmt.Printf("ERR!!\n\n")
		fmt.Println(message, err)
	}
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
