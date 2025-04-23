package mailstyler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/phzeng0726/gomailstyler/internal/service"
)

type MailManager interface {
	SendMail(mm MailMessage) error
	RenderTemplate(tmplFile string, data any) (string, error)
	RenderTemplateWithFuncs(tmplFile string, data any) (string, error)
	RenderTemplateWithCSS(tmplFile, cssFile string, data any) (string, error)
	RenderTemplateWithFuncsAndCSS(tmplFile, cssFile string, data any) (string, error)
}

type Manager struct {
	smtpServer string
	smtpPort   string
	smtpSender string
	tmplSvc    service.Templates
	cssToolSvc service.CSSTools
}

func NewManager(smtpServer, smtpPort, smtpSender, templatePath, cssPath string) (*Manager, error) {
	if smtpServer == "" {
		return nil, errors.New("empty smtpServer")
	}

	if smtpPort == "" {
		return nil, errors.New("empty smtpPort")
	}

	if smtpSender == "" {
		return nil, errors.New("empty smtpSender")
	}

	svc := service.NewServices(service.Deps{
		TmplPath: templatePath, // The root path of the template folder
		CSSPath:  cssPath,      // The root path of the css folder
	})

	return &Manager{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		smtpSender: smtpSender,
		tmplSvc:    svc.Templates,
		cssToolSvc: svc.CSSTools,
	}, nil
}

func (m *Manager) buildHTMLMessage(mm MailMessage) []byte {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	headers := map[string]string{
		"MIME-Version": "1.0",
		"From":         m.smtpSender,
		"To":           strings.Join(mm.To, ","),
		"Subject":      mm.Subject,
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s\n", boundary),
	}

	if len(mm.Cc) > 0 {
		headers["Cc"] = strings.Join(mm.Cc, ",")
	}

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// Write HTML message body
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("\r\n" + mm.Message)

	if len(mm.Attachments) > 0 {
		// Add each attachment as a separate MIME part
		for _, attachment := range mm.Attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", http.DetectContentType(attachment.Data)))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.FileName))
			buf.WriteString("\r\n")

			encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
			base64.StdEncoding.Encode(encoded, attachment.Data)

			// Encode attachment data as base64 and split into 76-character lines
			for i := 0; i < len(encoded); i += 76 {
				end := i + 76
				if end > len(encoded) {
					end = len(encoded)
				}
				buf.Write(encoded[i:end])
				buf.WriteString("\r\n")
			}
		}
	}

	// Final boundary to indicate end of MIME message
	buf.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	return buf.Bytes()
}

func (m *Manager) SendMail(mm MailMessage) error {
	addr := fmt.Sprintf("%s:%s", m.smtpServer, m.smtpPort)

	// Build HTML message
	msg := m.buildHTMLMessage(mm)

	// Sending email.
	allRecipients := append(mm.To, mm.Cc...)
	if err := smtp.SendMail(addr, nil, m.smtpSender, allRecipients, msg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) RenderTemplate(tmplFile string, data any) (string, error) {
	return m.tmplSvc.RenderTemplate(tmplFile, data)
}

func (m *Manager) RenderTemplateWithFuncs(tmplFile string, data any) (string, error) {
	return m.tmplSvc.RenderTemplateWithFuncs(tmplFile, data)
}

func (m *Manager) RenderTemplateWithCSS(tmplFile, cssFile string, data any) (string, error) {
	return m.cssToolSvc.RenderTemplateWithCSS(tmplFile, cssFile, data)
}

func (m *Manager) RenderTemplateWithFuncsAndCSS(tmplFile, cssFile string, data any) (string, error) {
	return m.cssToolSvc.RenderTemplateWithFuncsAndCSS(tmplFile, cssFile, data)

}
