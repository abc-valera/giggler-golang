package email

type emailer interface {
	sendEmail(e EmailSendReq) error
}

type EmailSendReq struct {
	Subject     string
	Content     string
	To          []string
	AttachFiles []string
}
