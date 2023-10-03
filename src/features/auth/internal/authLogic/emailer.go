package authLogic

import (
	"time"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/shared/logger"
)

var emailer = func() iEmailer {
	switch env.Load("EMAILER") {
	case "dummy":
		return infraDummyEmailer{}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}()

type (
	iEmailer interface {
		SendEmail(e emailSendReq) error
	}

	emailSendReq struct {
		Subject     string
		Content     string
		To          []string
		AttachFiles []string
	}
)

type infraDummyEmailer struct{}

func (infraDummyEmailer) SendEmail(e emailSendReq) error {
	// Sending email
	time.Sleep(250 * time.Millisecond)

	logger.Info("EMAIL_SENT", "to", e.To, "subject", e.Subject, "body", e.Content)
	return nil
}
