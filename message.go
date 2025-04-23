package mailstyler

type MailMessage struct {
	Subject      string
	Message      string
	To           []string
	Cc           []string
	Attachments  []Attachment
	InlineImages []InlineImage
}

type Attachment struct {
	FileName string
	Data     []byte
}

type InlineImage struct {
	CID      string
	FileName string
	Data     []byte
}
