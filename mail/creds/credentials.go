package creds

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
	"sort"
	"strings"
)

var dirsToCheck = []string{"./credentials", "../credentials", "../../credentials",
	"../../../credentials",
	"/credentials", "/etc/credentials"}

func walk(root string) (string, error) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "creds.json") {
			if strings.Contains(path, "token.json") {
				return nil
			} else {
				log.Fatalf("Run quickstart.go to" +
					" generate token.js" +
					"go run quickstart.go")
			}

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

type CREDS struct {
	b      []byte
	file   string
	dir    string
	client *http.Client
	token  []*oauth2.Token
	srv    *gmail.Service
}

func NewGmailSrv() *gmail.Service {
	c := CREDS{}
	c.PopulateCREDS()
	srv := c.GetSRV()
	return srv
}

func ListTokens(root string) ([]string, error) {
	files := []string{}
	token_files := []string{}


	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files,err
	}
	for _, file := range files {
		if strings.Contains(file,"token") {
			token_files = append(token_files,file)
		}
	}
	sort.Strings(token_files)
	return token_files, err
}



func (c *CREDS) PopulateCREDS() {
	dir, err := FindDir()
	if err != nil {
		log.Fatalf("Can't find credential file")
	}
	c.b = ReadCredentials(dir)

	token, err := tokenFromFile(dir + "/token.json")
	if err != nil {
		log.Printf("Can't read token.json. %v\n", err)
		log.Printf("Error msg: %v\n", err)
		token, err = tokenFromFile("/credentials/token.json")
		if err != nil {
			log.Printf("NOPE. NOT in /credentials/token.json")
			log.Printf("Error msg: %v\n", err)
			return
		}

		log.Printf("GOT IT.  /credentials/token.json")
	}
	c.token[0] = token
}


func (c *CREDS) PopulateCREDS2() {
	dir, err := FindDir()
	if err != nil {
		log.Fatalf("Can't find credential file")
	}
	c.b = ReadCredentials(dir)

	token, err := tokenFromFile(dir + "/token2.json")
	if err != nil {
		log.Printf("Can't read token.json. %v\n", err)
		log.Printf("Error msg: %v\n", err)
		token, err = tokenFromFile("/credentials/2token.json")
		if err != nil {
			log.Printf("NOPE. NOT in /credentials/2token.json")
			log.Printf("Error msg: %v\n", err)
			return
		}

		log.Printf("GOT IT.  /credentials/token2.json")
	}
	c.token[0] = token
}


func ReadCredentials(dir string) []byte {
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

func (c *CREDS) GetSRV() *gmail.Service {

	config, err := google.ConfigFromJSON(c.b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	c.client = config.Client(context.Background(), c.token[0])

	srv, err := gmail.New(c.client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	c.srv = srv
	return srv
}
