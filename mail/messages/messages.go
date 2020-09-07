package messages

import (
	"encoding/base64"
	"fmt"
	"github.com/mchirico/go-gmail/mail/creds"
	"google.golang.org/api/gmail/v1"
	"io"
	"strings"
)

// Labels - map of labels
func Labels() (map[string]string, error) {

	srv := creds.NewGmailSrv()
	user := "me"
	m := map[string]string{}
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		return m, err
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return m, err
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		m[l.Name] = l.Id
		fmt.Printf("- %s, %s\n", l.Name, l.Id)
	}
	return m, nil
}

func GetNewMessages(labelID string, maxCount int) []map[string]string {

	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)
	msg, _ := nsrv.Messages.List("me").LabelIds(labelID).Do()

	count := 0
	total := []map[string]string{}

	for _, v := range msg.Messages {

		count += 1
		if count > maxCount {
			break
		}

		g, _ := srv.Users.Messages.Get("me", v.Id).Format("metadata").Do()
		tag := map[string]string{}

		for _, v := range g.Payload.Headers {
			tag[v.Name] = v.Value
		}

		tag["Snippet"] = g.Snippet
		total = append(total, tag)
		tag["Id"] = v.Id
		total = append(total, tag)
	}
	return total
}

func GetRaw(labelID string, maxCount int) map[string][]byte {

	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)
	msg, _ := nsrv.Messages.List("me").LabelIds(labelID).Do()

	count := 0
	rmsg := map[string][]byte{}
	for _, v := range msg.Messages {
		count += 1
		if count > maxCount {
			break
		}

		g, _ := srv.Users.Messages.Get("me", v.Id).Format("raw").Do()
		data, _ := base64.RawURLEncoding.DecodeString(g.Raw)
		rmsg[v.Id] = data

	}
	return rmsg
}

type Message struct {
	From       string // from name
	ReplyTo    string // from address
	To         string
	Subject    string
	Body       string
	Attachment io.Reader
}

func Send(to string, subject string, body string) error {
	m := Message{}
	srv := creds.NewGmailSrv()
	var msg gmail.Message

	m.Subject = subject
	m.To = to
	m.Body = body
	s := "From: " + m.From + "\r\n" +
		"reply-to: " + m.ReplyTo + "\r\n" +
		"To: " + m.To + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		"\r\n" + m.Body

	msg.Raw = base64.StdEncoding.EncodeToString([]byte(s))
	_, err := srv.Users.Messages.Send("me", &msg).Do()
	return err

}

func Reply(replyID, msgID, from, to, subject, msg_to_send string) (string, error) {

	srv := creds.NewGmailSrv()

	rawMessage := ""
	rawMessage += fmt.Sprintf("To: %s\r\n", to)
	rawMessage += fmt.Sprintf("Subject: %s\r\n", subject)
	rawMessage += fmt.Sprintf("Reply-To: %s\r\n", from)
	rawMessage += fmt.Sprintf("In-Reply-To: %s\r\n", msgID)
	rawMessage += fmt.Sprintf("References: %s\r\n", msgID)
	rawMessage += fmt.Sprintf("Return-Path: %s\r\n", from)
	rawMessage += fmt.Sprintf("AI-Msg-Field: %s\r\n", "suspect")

	// Add extra linebreak for splitting headers and body
	rawMessage += "\r\n\r\n"
	rawMessage += msg_to_send

	// New message for our gmail service to send
	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString([]byte(rawMessage))
	message.ThreadId = replyID

	// Send the message
	_, err := srv.Users.Messages.Send("me", &message).Do()

	if err != nil {
		return "", err
	} else {
		return "Message sent!", err
	}

}

func Thread(labelID string, maxCount int) map[string][]byte {

	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)
	msg, _ := nsrv.Threads.List("me").LabelIds(labelID).Do()

	count := 0
	rmsg := map[string][]byte{}
	for _, v := range msg.Threads {
		count += 1
		if count > maxCount {
			break
		}
		fmt.Println(v.Id, v.HistoryId, v.Snippet)
		//g, _ := srv.Users.Messages.Get("me", v.Id).Format("raw").Do()
		//data, _ := base64.RawURLEncoding.DecodeString(g.Raw)
		//rmsg[v.Id] = data

	}
	return rmsg
}

func Domains(r []map[string]string) map[string]int {
	domains := map[string]int{}
	for id := range r {
		s := r[id]["Return-Path"]
		if strings.ContainsAny(s, "0123456789-") {
			continue
		}
		idx0 := strings.Index(s, "@")
		idx1 := strings.Index(s, ">")
		if idx0 > 1 && idx1 > 1 {
			domains[s[idx0+1:idx1]] += 1
		}

	}
	return domains
}

func Watch(userId string, watchReq *gmail.WatchRequest) (*gmail.UsersWatchCall) {

	srv := creds.NewGmailSrv()
    nsrv := gmail.NewUsersService(srv)
	return nsrv.Watch(userId, watchReq)

}

func StopWatch(userId string) (error) {
	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)
	return nsrv.Stop(userId).Do()
}

