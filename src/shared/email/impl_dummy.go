package email

import (
	"time"
)

type dummy struct{}

func (dummy) sendEmail(e EmailSendReq) error {
	// Sending email...
	time.Sleep(250 * time.Millisecond)

	return nil
}
