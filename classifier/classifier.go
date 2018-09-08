package classifier

import (
	. "Mail/definitions"
	"github.com/emersion/go-imap"
)

// ------------------------------------------- Public ------------------------------------------- //

// Given email and possible categories, returns urgency level and category of email's contents
func Classify(email *imap.Message, categories []Category) <-chan Classification {
	returnChan := make(chan Classification)

	go func() {
		defer close(returnChan)
		//
		// classification := some python code...
		//
		returnChan <- Classification{Email: email, Urgent: URGENT, Type: 0}
	}()

	return returnChan
}
