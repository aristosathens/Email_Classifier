package definitions

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// ------------------------------------------- Types ------------------------------------------- //

// User defined categories (soccer, class, fraternity, etc.)
type Category int

// Level of response needed
type Urgency int

const (
	NORMAL Urgency = 0
	URGENT Urgency = 1
	OTHER  Urgency = 2
)

type Classification struct {
	Email  *imap.Message
	Urgent Urgency
	Type   Category
}

// ------------------------------------------- User ------------------------------------------- //

type User struct {
	client     *client.Client
	Emails     []*imap.Message
	inbox      *imap.MailboxStatus
	mailboxes  chan *imap.MailboxInfo
	Categories []Category
}

// Sets up client object
func (c *User) GetClient(user, password string) {

	// Connect to server
	cl, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	c.client = cl

	// Login
	if err := c.client.Login(user, password); err != nil {
		log.Fatal(err)
	}
}

func (c *User) GetInbox() {

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.client.List("", "*", mailboxes)
	}()

	// Ensure graceful exit
	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select inbox
	in, err := c.client.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	c.inbox = in
}

// Gets all emails from c.inbox, stores in c.emails
func (c *User) GetEmails() {

	seqset := new(imap.SeqSet)
	seqset.AddRange(1, c.inbox.Messages)
	items := []imap.FetchItem{imap.FetchEnvelope}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.client.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		c.Emails = append(c.Emails, msg)
	}
	if err := <-done; err != nil {
		log.Fatal(err)
	}
}
