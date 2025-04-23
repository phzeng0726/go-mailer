package mailstyler

type MailMessage struct {
	Subject string
	Message string
	To      []string
	Cc      []string
}
