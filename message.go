package mailstyler

type MailMessage struct {
	Subject     string
	Message     string
	To          []string
	Cc          []string
	Attachments []Attachment
}

type Attachment struct {
	FileName string
	Data     []byte
}
