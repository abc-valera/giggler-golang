package emailer

import (
	"fmt"
	"time"
)

type dummy struct{}

func InitDummy() Interface {
	return dummy{}
}

func (dummy) Send(e EmailData) error {
	time.Sleep(250 * time.Millisecond)

	fmt.Println("Email Sent")
	fmt.Println("\tTo:", e.To)
	fmt.Println("\tSubject:", e.Subject)
	fmt.Println("\tContent:", e.Content)

	return nil
}
