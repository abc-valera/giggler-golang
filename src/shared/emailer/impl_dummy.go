package emailer

import (
	"fmt"
	"time"
)

type dummy struct{}

func (dummy) Send(e EmailSendIn) error {
	// Sending email...
	time.Sleep(250 * time.Millisecond)

	fmt.Println("Email Sent")
	fmt.Println("To: ", e.To)
	fmt.Println("Subject: ", e.Subject)
	fmt.Println("Content: ", e.Content)

	return nil
}
