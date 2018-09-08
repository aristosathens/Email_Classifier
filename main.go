//
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"Mail/classifier"
	. "Mail/definitions"
)

// ------------------------------------------- Main ------------------------------------------- //

func main() {

	// Use dummy account: Aristos.Website, VerySecurePassword
	c := User{}
	userName := getUserInput("Enter email: ")
	password := getUserInput("Enter password: ")
	c.GetClient(userName, password)
	c.GetInbox()
	c.GetEmails()

	// Send every email to classifier. Will get channel for each in response
	channels := make([]<-chan Classification, len(c.Emails))
	fmt.Println("Unread emails: ")
	for i, email := range c.Emails {
		channels[i] = classifier.Classify(email, c.Categories)
		fmt.Println("Subject: " + email.Envelope.Subject)
	}

	// Merge all channels so we only have to listen on a single channel
	ch := mergeChannels(channels)

	// Listen on our channel until all emails are classified
	for {
		select {
		case classifiedEmail, ok := <-ch:
			//
			// send to front end...
			//
			fmt.Println(classifiedEmail)
			if !ok {
				break
			}
		}
	}
}

// ------------------------------------------- Private ------------------------------------------- //

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
