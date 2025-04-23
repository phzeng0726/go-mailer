# Go Mail Styler

Go Mail Styler is a simple and flexible Go package for sending HTML emails via SMTP.

It provides an easy way to render HTML templates with dynamic data and send them as emails.

---

## ✨ Features

- Simple API for sending emails
- HTML template rendering with dynamic data
- Support for built-in and custom template functions
- Easily apply a CSS file across all templates
- Embed attachments and images directly within HTML
- Lightweight and easy to integrate

## 📦 Installation

```bash
go get github.com/phzeng0726/gomailstyler@v0.2.1
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

	// Optional: Load image as []byte
	imageData, err := fileToBytes("./assets/images/my_doggy.jpg")
	if err != nil {
		log.Fatalf("failed to load image: %v", err)
	}

	// Send the email
	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Hello",
		Message: body,
		To:      []string{"someone@example.com"},
		Cc:      []string{"someone@example.com"}, // Optional
		Attachments: []mailstyler.Attachment{  // Optional
			{
				FileName: "my_doggy.jpg",
				Data:     imageData,
			},
		},
		InlineImages: []mailstyler.InlineImage{  // Optional
			{
				CID:      "my-doggy-img",
				FileName: "my_doggy.jpg",
				Data:     imageData,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}

	log.Println("Mail sent successfully!")
}
```

---

## 🏗️ Input Struct

### MailMessage

A `MailMessage` represents the email that will be sent.

```go
type MailMessage struct {
    Subject      string          // The subject of the email
    Message      string          // The body content of the email in HTML format
    To           []string        // List of recipient email addresses
    Cc           []string        // List of CC email addresses (optional)
    Attachments  []Attachment    // List of file attachments (optional)
    InlineImages []InlineImage   // List of inline images to be embedded (optional)
}
```

### Attachment

An `Attachment` represents a file to be attached to the email.

```go
type Attachment struct {
    FileName string   // The name of the attachment file
    Data     []byte   // The binary data of the attachment file
}
```

### InlineImage

An `InlineImage` represents an image that is embedded within the HTML body of the email.

```go
type InlineImage struct {
    CID      string   // The content ID to reference the image within the HTML (e.g., cid:image-id)
    FileName string   // The name of the image file
    Data     []byte   // The binary data of the image
}

```

---

## 📘 API Reference

### `NewManager`

```go
func NewManager(
	smtpHost, smtpPort, smtpSender string,
	smtpPassword *string,
	templatePath, cssPath string,
) (*Manager, error)
```

Creates a new mail manager instance.

- `smtpHost`: SMTP server host address (e.g., `smtp.example.com`).

  - This specifies the address of the SMTP server used for sending emails.

- `smtpPort`: SMTP server port (e.g., `587`).

  - This specifies the port number used to connect to the SMTP server. Typically, `587` is used for secure email transmission with STARTTLS, and `465` is for SSL/TLS.

- `smtpSender`: Sender's email address.

  - The email address that will be used as the "From" address when sending emails.

- `smtpPassword`: Sender's SMTP password (or application-specific password if applicable).

  - This is the password (or application-specific password for Gmail and similar services) required to authenticate with the SMTP server.

- `templatePath`: Root directory of your template files.

  - The file system path where email templates are stored. These templates will be used to render email content dynamically.

- `cssPath`: Root directory of your CSS files.
  - The file system path where CSS files are stored. These CSS files will be used to style email templates when sending HTML emails.

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
func (m *Manager) RenderTemplateWithFuncs(tmplFile string, data any, customFuncs ...template.FuncMap) (string, error)
```

Renders an HTML template with additional template functions. You can define your own custom functions by passing a template.FuncMap as customFuncs. If no customFuncs are provided, the package will use a set of default, commonly used functions.

---

### `RenderTemplateWithCSS`

```go
func (m *Manager) RenderTemplateWithCSS(tmplFile, cssFile string, data any) (string, error)
```

Renders an HTML template using the provided data and an external CSS file, automatically converting the CSS to inline styles.

---

### `RenderTemplateWithFuncsAndCSS`

```go
func (m *Manager) RenderTemplateWithFuncsAndCSS(tmplFile, cssFile string, data any, customFuncs ...template.FuncMap) (string, error)
```

Renders an HTML template with additional template functions and an external CSS file, automatically converting the CSS to inline styles.
Similar to `RenderTemplateWithFuncs`, you can provide your own custom functions, and the package will fall back to the default ones if none are provided.

---

### 🧰 Default Functions

| Function     | Description                                                     | Example                                            |
| ------------ | --------------------------------------------------------------- | -------------------------------------------------- |
| `add`        | Adds two integers                                               | `{{ add 1 2 }}` → `3`                              |
| `toUpper`    | Converts a string to uppercase                                  | `{{ toUpper "hello" }}` → `HELLO`                  |
| `toLower`    | Converts a string to lowercase                                  | `{{ toLower "hello" }}` → `hello`                  |
| `trim`       | Removes leading and trailing whitespace                         | `{{ trim " hello " }}` → `hello`                   |
| `title`      | Capitalizes the first letter of each word                       | `{{ title "go mail styler" }}` → `Go Mail Styler`  |
| `formatDate` | Formats a `time.Time` using a Go layout string                  | `{{ formatDate now "2006-01-02" }}` → `2025-04-23` |
| `isEmpty`    | Returns `true` if a string is empty or contains only whitespace | `{{ isEmpty "  " }}` → `true`                      |
| `safeHTML`   | Marks a string as safe HTML (not escaped in output)             | `{{ safeHTML "<b>bold</b>" }}`                     |
| `inc`        | Increments an integer by 1                                      | `{{ inc 0 }}` → `1`                                |

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
