package emailer

import (
	"time"
)

type dummy struct{}

func (dummy) Send(e EmailSendIn) error {
	// Sending email...
	time.Sleep(250 * time.Millisecond)

	return nil
}
