package emailer

type Interface interface {
	Send(e SendIn) error
}

type SendIn struct {
	To          []string
	Subject     string
	Content     string
	AttachFiles []string
}
