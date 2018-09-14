package classifier

import (
	. "Mail/definitions"
	// "encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"os/exec"
)

// ------------------------------------------- Types ------------------------------------------- //

type SendPackage struct {
	Date    string
	Senders []string
	Subject string
	Body    string
}

// ------------------------------------------- Public ------------------------------------------- //

// Given email and possible categories, returns urgency level and category of email's contents
func Classify(email *imap.Message, categories []Category) <-chan Classification {
	returnChan := make(chan Classification)

	go func() {

		defer close(returnChan)

		//
		// classification := some python code...
		//
		exec.Command("classifier_script.py").Run()

		returnChan <- Classification{Email: email, Urgent: URGENT, Type: 0}
	}()

	return returnChan
}

func main(email *imap.Message) {

	// email.Envelope.
	byteArray := []byte("hi there!")

	cmd := exec.Command("python3", "mail_classifier.py")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		if _, err := stdin.Write(byteArray); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Exec status: ", cmd.Run())
}

// func emailToBytes(email *imap.Message) []byte {

// 	senders := make([]string, len(email.Envelope.Sender))
// 	for i, s := range email.Envelope.Sender {
// 		senders[i] = s.PersonalName
// 	}

// 	pkg := SendPackage{
// 		Date:    "1/1/1970",
// 		Senders: senders,
// 		Subject: email.Envelope.Subject,
// 		Body:    email.Body,
// 	}

// 	// senders := ""
// 	for _, s := range email.Envelope.Sender {
// 		senders += s.PersonalName + ","
// 	}

// 	sender := []byte(email.Envelope.Sender)
// 	subject := []byte("Hello Aristos")
// 	body := []byte("This is an important email")
// 	toSend := [][]byte{bDelimiter, sender, bDelimiter, subject, bDelimiter, body}
// 	arr := bDelimiter
// 	arr = bytes.Join(toSend, arr)
// 	fmt.Println(string(arr))
// }

// func sendersToString(senders *[]*Address) *[]string {
// 	for _, s := range senders {

// 	}
// }
