package mail

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var dirsToCheck = []string{"./credentials", "../credentials",
	"/credentials", "/etc/credentials"}

func walk(root string) (string, error) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "credentials.json") {
			return nil
		}
		return err
	})
	return root, err
}

func FindDir() (string, error) {

	for _, dir := range dirsToCheck {
		files, err := walk(dir)
		if err == nil {
			return files, err
		}
	}
	return "", errors.New("not found")
}

func ReadCredentials() []byte {
	dir, err := FindDir()
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	b, err := ioutil.ReadFile(dir + "/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	return b
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getClient(config *oauth2.Config, dir string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := dir + "/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Fatalf("getClient: %v", err.Error())
	}
	return config.Client(context.Background(), tok)
}

func GetSRV() (*gmail.Service) {

	config, err := google.ConfigFromJSON(ReadCredentials(), gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	dir, err := FindDir()
	client := getClient(config, dir)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv
}
