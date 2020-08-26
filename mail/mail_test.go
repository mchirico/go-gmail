package mail

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"log"
	"testing"
)

func Test_findDir(t *testing.T) {

	dir, err := FindDir()
	fmt.Println(dir)
	if err != nil {
		t.Errorf("Can't find .json")
	}

}

func Test_GetSRV(t *testing.T) {

	srv := GetSRV()
	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
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
	msg, _ := srv.Users.Messages.List("me").LabelIds("TRASH").Do()

	nsrv := gmail.NewUsersService(srv)
	nsrv.Messages.List("me").LabelIds("TRASH").Do()

	//msg, _ := srv.Users.Messages.List("me").Do()

	count := 0
	for _, v := range msg.Messages {
		//fmt.Println(v)
		count += 1
		if count > 9 {
			continue
		}
		g, _ := srv.Users.Messages.Get("me", v.Id).Format("raw").Do()
		fmt.Println(g.Snippet)
		data, _ := base64.StdEncoding.DecodeString(g.Raw)
		sdata := fmt.Sprintf("%s",data)

		fmt.Printf("\n%s\n\n____\n",sdata[0:])


	}

}
