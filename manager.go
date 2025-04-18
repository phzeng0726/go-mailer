package mailstyler

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"

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

func (m *Manager) buildHTMLMessage(subject, body string) []byte {
	headers := map[string]string{
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
		"Subject":      subject,
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	msg.WriteString("\r\n" + body)

	return msg.Bytes()
}

func (m *Manager) SendMail(mm MailMessage) error {
	addr := fmt.Sprintf("%s:%s", m.smtpServer, m.smtpPort)

	// Connect to SMTP server
	msg := m.buildHTMLMessage(mm.Subject, mm.Message)

	// Sending email.
	if err := smtp.SendMail(addr, nil, m.smtpSender, mm.To, msg); err != nil {
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
