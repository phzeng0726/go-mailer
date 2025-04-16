package mail

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"path/filepath"
	"text/template"
)

type MailManager interface {
	SendMail(mm MailMessage) error
	RenderTemplate(templateFile string, data any) (string, error)
	RenderTemplateWithFuncs(templateFile string, data any) (string, error)
}

type Manager struct {
	smtpServer   string
	smtpPort     string
	smtpSender   string
	templatePath string // The root path of the template folder
}

type MailMessage struct {
	Subject  string
	Message  string
	MsgStyle string
	To       []string
}

func NewManager(smtpServer string, smtpPort string, smtpSender string, templatePath string) (*Manager, error) {
	if smtpServer == "" {
		return nil, errors.New("empty smtpServer")
	}

	if smtpPort == "" {
		return nil, errors.New("empty smtpPort")
	}

	if smtpSender == "" {
		return nil, errors.New("empty smtpSender")
	}

	return &Manager{
		smtpServer:   smtpServer,
		smtpPort:     smtpPort,
		smtpSender:   smtpSender,
		templatePath: templatePath,
	}, nil
}

func (m *Manager) SendMail(mm MailMessage) error {
	addr := fmt.Sprintf("%s:%s", m.smtpServer, m.smtpPort)

	msgStr := fmt.Sprintf("Subject: %s\r\n%s%s", mm.Subject, mm.MsgStyle, mm.Message)

	// Connect to SMTP server
	msg := []byte(msgStr)

	// Sending email.
	if err := smtp.SendMail(addr, nil, m.smtpSender, mm.To, msg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) RenderTemplate(templateFile string, data any) (string, error) {
	tmplPath := filepath.Join(m.templatePath, templateFile)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

func (m *Manager) RenderTemplateWithFuncs(templateFile string, data any) (string, error) {
	tmplPath := filepath.Join(m.templatePath, templateFile)

	// Register the "add" function in the template so that the index can start from 1
	tmpl := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	tmplWithAdd, err := tmpl.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmplWithAdd.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
