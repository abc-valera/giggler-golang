package auth

import (
	"time"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/shared/logger"
)

var emailer iEmailer

func init() {
	switch env.Load("EMAILER") {
	case "dummy":
		emailer = infraDummyEmailer{}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type (
	iEmailer interface {
		SendEmail(e emailSendRequest) error
	}

	emailSendRequest struct {
		Subject     string
		Content     string
		To          []string
		AttachFiles []string
	}
)

type infraDummyEmailer struct{}

func (infraDummyEmailer) SendEmail(e emailSendRequest) error {
	// Sending email
	time.Sleep(250 * time.Millisecond)

	logger.Info("EMAIL_SENT", "to", e.To, "subject", e.Subject, "body", e.Content)
	return nil
}
