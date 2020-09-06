package messages

import (
	"encoding/base64"
	"fmt"
	"github.com/mchirico/go-gmail/mail/creds"
	"google.golang.org/api/gmail/v1"
	"io"
	"log"
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

func GetNewMessages(labelID string, maxCount int) []map[string][]string {

	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)
	msg, _ := nsrv.Messages.List("me").LabelIds(labelID).Do()

	count := 0
	total := []map[string][]string{}

	for _, v := range msg.Messages {

		count += 1
		if count > maxCount {
			break
		}

		g, _ := srv.Users.Messages.Get("me", v.Id).Format("metadata").Do()
		//fmt.Println(g.Snippet)

		tag := map[string][]string{}

		for _, v := range g.Payload.Headers {
			if r, ok := tag[v.Name]; ok {
				r = append(r, v.Value)
			} else {
				tag[v.Name] = []string{v.Value}
			}
		}

		tag["Snippet"] = []string{g.Snippet}
		total = append(total, tag)
		tag["Id"] = []string{v.Id}
		total = append(total,tag)
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

func Reply(replyID, from, to, msg_to_send string)  {

	srv := creds.NewGmailSrv()
	nsrv := gmail.NewUsersService(srv)


	// replyID := "174652e40183c2e9"

	msg, _ := nsrv.Messages.Get("me", replyID).Format("metadata").Do()


	var rawMessage = []string{}

	//from := "dead@cwxstat.com"
	//to := "mchirico@gmail.com"

	subject := ""
	msgID := ""
	for _, v := range msg.Payload.Headers {
		if v.Name == "Subject" {
			subject = v.Value
		}
		if v.Name == "Message-ID" {
			msgID = v.Value
		}
	}

	subject = "C2C Contracts Only...  Re: " + subject

	// Add the to and Reply-To
	rawMessage = append(rawMessage, fmt.Sprintf("To: %s\r\n", to))
	rawMessage = append(rawMessage, fmt.Sprintf("Subject: %s\r\n", subject))
	rawMessage = append(rawMessage, fmt.Sprintf("Reply-To: %s\r\n", from))
	rawMessage = append(rawMessage, fmt.Sprintf("In-Reply-To: %s\r\n", msgID))
	rawMessage = append(rawMessage, fmt.Sprintf("References: %s\r\n", msgID))
	rawMessage = append(rawMessage, fmt.Sprintf("Return-Path: %s\r\n", from))


	// Add extra linebreak for splitting headers and body
	rawMessage = append(rawMessage, "\r\n\r\n")

	rawMessage = append(rawMessage, msg_to_send)

	// New message for our gmail service to send
	var message gmail.Message
	messageStr := []byte(strings.Join(rawMessage, ""))
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)
	message.ThreadId = replyID

	// Send the message
	_, err := srv.Users.Messages.Send("me", &message).Do()

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}

}


