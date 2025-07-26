package email

type emailer interface {
	sendEmail(e EmailSendIn) error
}

type EmailSendIn struct {
	Subject     string
	Content     string
	To          []string
	AttachFiles []string
}
