# Go Mail Styler

Go Mail Styler is a simple and flexible Go package for sending HTML emails via SMTP.

It provides an easy way to render HTML templates with dynamic data and send them as emails.

## 📦 Installation

```bash
go get github.com/phzeng0726/gomailstyler@v0.1.5
```

## 🚀 Example Usage

```go
package main

import (
	"log"
	"github.com/phzeng0726/gomailstyler"
)

func main() {
	// Initialize Mail Manager
	manager, err := mailstyler.NewManager("smtp.example.com", "587", "you@example.com", "./templates", "./templates/css")
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}

	// Render template with dynamic data
	body, err := manager.RenderTemplate("welcome.html", map[string]any{
		"Name": "Someone",
	})
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}

	// Send the email
	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Hello",
		Message: body,
		To:      []string{"someone@example.com"},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}

	log.Println("Mail sent successfully!")
}
```

---

## ✨ Features

- Simple API for sending emails
- HTML template rendering with dynamic data
- Support for template helper functions
- Lightweight and easy to integrate

---

## 📘 API Reference

### `NewManager`

```go
func NewManager(smtpServer, smtpPort, smtpSender, templatePath, cssPath string) (*Manager, error)
```

Creates a new mail manager instance.

- `smtpServer`: SMTP server address (e.g., `smtp.example.com`)
- `smtpPort`: SMTP server port (e.g., `587`)
- `smtpSender`: Sender's email address
- `templatePath`: Root directory of your template files
- `cssPath`: Root directory of your css files

---

### `SendMail`

```go
func (m *Manager) SendMail(mm MailMessage) error
```

Sends an email.

- `mm`: A `MailMessage` struct containing subject, message body, and recipient addresses.

---

### `RenderTemplate`

```go
func (m *Manager) RenderTemplate(tmplFile string, data any) (string, error)
```

Renders an HTML template using the provided data.

---

### `RenderTemplateWithFuncs`

```go
func (m *Manager) RenderTemplateWithFuncs(tmplFile string, data any) (string, error)
```

Renders an HTML template with additional template functions like `add`.

---

### `RenderTemplateWithCSS`

```go
func (m *Manager) RenderTemplateWithCSS(tmplFile, cssFile string, data any) (string, error)
```

Renders an HTML template using the provided data and an external CSS file, automatically converting the CSS to inline styles.

---

### `RenderTemplateWithFuncsAndCSS`

```go
func (m *Manager) RenderTemplateWithFuncsAndCSS(tmplFile, cssFile string, data any) (string, error)
```

Renders an HTML template with additional template functions like `add` and an external CSS file, automatically converting the CSS to inline styles.

---

## ✉️ MailMessage Struct

```go
type MailMessage struct {
	Subject string   // Email subject
	Message string   // HTML email content
	To      []string // Recipient email addresses
}
```

---

## 📝 Example Template (`hello.html`)

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Hello, {{.Name}}</title>
  </head>
  <body>
    <h1>Hello, {{.Name}}!</h1>
    <p>Welcome to our service.</p>
  </body>
</html>
```

This template will be rendered with the `Name` field provided in the `data` map.

---

## 🪪 License

[MIT](./LICENSE)
