package email

import (
	"giggler-golang/src/shared/errutil/must"
)

var emailerInstance = func() emailer {
	switch must.Env("EMAILER") {
	case "dummy":
		return dummy{}
	default:
		panic(must.ErrInvalidEnvValue)
	}
}()

func Send(e EmailSendIn) error {
	return emailerInstance.sendEmail(e)
}
