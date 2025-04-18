package mail

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vanng822/go-premailer/premailer"
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
	Subject string
	Message string
	To      []string
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

func (m *Manager) RenderTemplate(templateFile, cssFile string, data any) (string, error) {
	tmplPath := filepath.Join(m.templatePath, templateFile)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// 讀本地 CSS 檔案
	cssPath := filepath.Join(m.templatePath, cssFile)
	cssBytes, err := os.ReadFile(cssPath)
	if err != nil {
		return "", fmt.Errorf("failed to read css file %s: %w", cssFile, err)
	}
	cssContent := string(cssBytes)

	// 插入 CSS 到 <head><style> 中（粗略做法：用字串替換）
	htmlWithCSS := strings.Replace(
		buf.String(),
		"</head>",
		fmt.Sprintf("<style>%s</style></head>", cssContent),
		1,
	)

	// Convert css into inline style
	options := premailer.NewOptions()
	options.RemoveClasses = false  // 不移除 class（可選）
	options.CssToAttributes = true // 將 CSS 屬性轉為 HTML 屬性（支援更廣）

	prem, err := premailer.NewPremailerFromString(htmlWithCSS, options)
	if err != nil {
		log.Fatal(err)
	}

	html, err := prem.Transform()
	if err != nil {
		log.Fatal(err)
	}

	return html, nil
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

	var buf bytes.Buffer
	if err := tmplWithAdd.Execute(&buf, data); err != nil {
		return "", err
	}

	// Convert css into inline style
	prem, err := premailer.NewPremailerFromString(buf.String(), premailer.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	html, err := prem.Transform()
	if err != nil {
		log.Fatal(err)
	}

	return html, nil
}
