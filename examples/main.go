package main

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	mailstyler "github.com/phzeng0726/gomailstyler"
)

var (
	cfg     *Config
	manager *mailstyler.Manager
)

// Config holds the environment variables for the mailer
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPSender   string
	SMTPPassword string // optional
	MailReceiver string
	TemplatePath string
	CSSPath      string
}

// loadConfig loads environment variables into a Config struct
func loadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),     // e.g., "smtp.example.com"
		SMTPPort:     os.Getenv("SMTP_PORT"),     // e.g., "587"
		SMTPSender:   os.Getenv("SMTP_SENDER"),   // e.g., "you@example.com"
		SMTPPassword: os.Getenv("SMTP_PASSWORD"), // e.g., "you@example.com"
		MailReceiver: os.Getenv("MAIL_RECEIVER"), // e.g., "recipient@example.com"
		TemplatePath: "./assets/templates",       // Root folder for templates
		CSSPath:      "./assets/templates/css",   // Root folder for CSS
	}
}

func fileToBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// Examples ==============================================
func sendBasicMail(manager *mailstyler.Manager, receiver string) {
	body, err := manager.RenderTemplate("welcome.html", map[string]any{
		"Name": "Alice",
		"Age":  30,
	})
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}

	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Rendered Template Email",
		Message: body,
		To:      []string{receiver},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}
}

func sendFuncsMail(manager *mailstyler.Manager, receiver string) {
	body, err := manager.RenderTemplateWithFuncs("welcome_with_funcs.html", map[string]any{
		"Name": "Bob",
	})
	if err != nil {
		log.Fatalf("failed to render template with funcs: %v", err)
	}

	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Template with Functions Email",
		Message: body,
		To:      []string{receiver},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}
	log.Println("Template with functions mail sent successfully!")
}

func sendCSSMail(manager *mailstyler.Manager, receiver string) {
	body, err := manager.RenderTemplateWithCSS("welcome.html", "styles.css", map[string]any{
		"Name": "Charlie",
	})
	if err != nil {
		log.Fatalf("failed to render template with CSS: %v", err)
	}

	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Template with CSS Email",
		Message: body,
		To:      []string{receiver},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}
	log.Println("Template with CSS mail sent successfully!")
}

func sendFuncsAndCSSMail(manager *mailstyler.Manager, receiver string) {
	body, err := manager.RenderTemplateWithFuncsAndCSS("welcome_with_funcs.html", "styles.css", map[string]any{
		"Name": "David",
	})
	if err != nil {
		log.Fatalf("failed to render template with funcs and CSS: %v", err)
	}

	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Template with Functions and CSS Email",
		Message: body,
		To:      []string{receiver},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}
	log.Println("Template with functions and CSS mail sent successfully!")
}

func sendCSSAndAttachmentsMail(manager *mailstyler.Manager, receiver string) {
	body, err := manager.RenderTemplateWithCSS("welcome_with_inline_images.html", "styles.css", map[string]any{
		"Name": "Pipi",
	})
	if err != nil {
		log.Fatalf("failed to render template with CSS: %v", err)
	}

	imageData, err := fileToBytes("./assets/images/my_doggy.jpg")
	if err != nil {
		log.Fatalf("failed to load image: %v", err)
	}

	err = manager.SendMail(mailstyler.MailMessage{
		Subject: "Template with CSS Email With Attachments",
		Message: body,
		To:      []string{receiver},
		Attachments: []mailstyler.Attachment{
			{
				FileName: "my_doggy.jpg",
				Data:     imageData,
			},
		},
		InlineImages: []mailstyler.InlineImage{
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
	log.Println("Template with CSS And Attachments mail sent successfully!")
}

// Run ==============================================
func init() {
	// Load configuration
	cfg = loadConfig()

	// Initialize mailer manager
	var err error
	manager, err = mailstyler.NewManager(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPSender,
		&cfg.SMTPPassword,
		cfg.TemplatePath,
		cfg.CSSPath,
	)
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}

}

func main() {
	sendBasicMail(manager, cfg.MailReceiver)
	sendFuncsMail(manager, cfg.MailReceiver)
	sendCSSMail(manager, cfg.MailReceiver)
	sendFuncsAndCSSMail(manager, cfg.MailReceiver)
	sendCSSAndAttachmentsMail(manager, cfg.MailReceiver)
}
