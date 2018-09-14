package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	. "Mail/definitions"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	msgs, err := srv.Users.Messages.List(user).Do()
	for _, msg := range msgs.Messages {
		if msg.Payload != nil {
			fmt.Println(msg.Payload.Body.Data)
		}
	}

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
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// //
// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"Mail/classifier"
// 	. "Mail/definitions"
// )

// // ------------------------------------------- Main ------------------------------------------- //

// func main() {

// 	// Use dummy account: Aristos.Website, VerySecurePassword
// 	c := User{}
// 	userName := getUserInput("Enter email: ")
// 	password := getUserInput("Enter password: ")
// 	c.GetClient(userName, password)
// 	c.GetInbox()
// 	c.GetEmails()

// 	// Send every email to classifier. Will get channel for each in response
// 	channels := make([]<-chan Classification, len(c.Emails))
// 	fmt.Println("Unread emails: ")
// 	for i, email := range c.Emails {
// 		channels[i] = classifier.Classify(email, c.Categories)
// 		fmt.Println("Subject: " + email.Envelope.Subject)
// 	}

// 	// Merge all channels so we only have to listen on a single channel
// 	ch := mergeChannels(channels)

// 	// Listen on our channel until all emails are classified
// 	for {
// 		select {
// 		case classifiedEmail, ok := <-ch:
// 			//
// 			// send to front end...
// 			//
// 			fmt.Println(classifiedEmail)
// 			if !ok {
// 				break
// 			}
// 		}
// 	}
// }

// // ------------------------------------------- Private ------------------------------------------- //

// Get user input
func getUserInput(message string) string {

	if message != "" {
		fmt.Print(message)
	}
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	input = strings.TrimSpace(input)
	return input
}

// Returns a channel that listens to all channels in the input array
func mergeChannels(channels []<-chan Classification) <-chan Classification {

	aggregateChannel := make(chan Classification, len(channels))
	for _, ch := range channels {
		go func(c <-chan Classification) {
			for {
				msg, flag := <-c
				if !flag {
					break
				}
				aggregateChannel <- msg
			}
		}(ch)
	}
	return aggregateChannel
}
