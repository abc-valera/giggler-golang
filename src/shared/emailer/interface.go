package emailer

type Interface interface {
	Send(e EmailData) error
}

type EmailData struct {
	To      []string
	Subject string
	Content string
}
