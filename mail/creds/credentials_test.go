package creds

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"log"
	"testing"
)

func TestListTokens(t *testing.T) {
	dir, _ := FindDir()
	r, _ := ListTokenFiles(dir)
	fmt.Printf("%v\n", r)
}

func Test_findDir(t *testing.T) {

	c := CREDS{}
	c.PopulateCREDS()
	if len(c.token) == 0 {
		t.Fatal()
	}

}

func Test_GetSRV(t *testing.T) {

	srv, _ := NewGmailSrv()

	user := "me"
	r, err := srv[0].Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
	msg, _ := srv[0].Users.Messages.List("me").LabelIds("TRASH").Do()

	nsrv := gmail.NewUsersService(srv[0])
	nsrv.Messages.List("me").LabelIds("TRASH").Do()

	//msg, _ := srv.Users.Messages.List("me").Do()

	count := 0
	for _, v := range msg.Messages {
		//fmt.Println(v)
		count += 1
		if count > 9 {
			continue
		}
		g, _ := srv[0].Users.Messages.Get("me", v.Id).Format("raw").Do()
		fmt.Println(g.Snippet)
		data, _ := base64.StdEncoding.DecodeString(g.Raw)
		sdata := fmt.Sprintf("%s", data)

		fmt.Printf("\n%s\n\n____\n", sdata[0:])

	}

}
