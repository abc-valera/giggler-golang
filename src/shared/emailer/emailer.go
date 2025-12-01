package emailer

import (
	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/singleton"
)

var Get = singleton.New(func() emailer {
	switch must.GetEnv("EMAILER") {
	case "dummy":
		return dummy{}
	default:
		panic(must.ErrInvalidEnvValue)
	}
})

type emailer interface {
	Send(e EmailSendIn) error
}

type EmailSendIn struct {
	To          []string
	Subject     string
	Content     string
	AttachFiles []string
}
