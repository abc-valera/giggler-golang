package email

import (
	"giggler-golang/src/shared/env"
)

var emailerInstance = func() emailer {
	switch env.Load("EMAILER") {
	case "dummy":
		return dummy{}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}()

func Send(e EmailSendIn) error {
	return emailerInstance.sendEmail(e)
}
