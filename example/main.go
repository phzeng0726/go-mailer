package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	mailstyler "github.com/phzeng0726/gomailstyler"
)

// Config holds the environment variables for the mailer
type Config struct {
	SMTPServer   string
	SMTPPort     string
	SMTPSender   string
	MailReceiver string
	TemplatePath string
	CSSPath      string
}

// loadConfig loads environment variables into a Config struct
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		SMTPServer:   os.Getenv("SMTP_SERVER"),   // e.g., "smtp.example.com"
		SMTPPort:     os.Getenv("SMTP_PORT"),     // e.g., "587"
		SMTPSender:   os.Getenv("SMTP_SENDER"),   // e.g., "you@example.com"
		MailReceiver: os.Getenv("MAIL_RECEIVER"), // e.g., "recipient@example.com"
		TemplatePath: "./templates",              // Root folder for templates
		CSSPath:      "./templates/css",          // Root folder for CSS
	}
}

// createManager initializes the mailstyler Manager
func createManager(cfg Config) *mailstyler.Manager {
	manager, err := mailstyler.NewManager(cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPSender, cfg.TemplatePath, cfg.CSSPath)
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}
	return manager
}

// sendTemplateMail demonstrates the RenderTemplate function
func sendTemplateMail(manager *mailstyler.Manager, receiver string) {
	log.Println("\nExample 2: Using RenderTemplate")
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
	log.Println("Template mail sent successfully!")
}

// sendTemplateWithFuncsMail demonstrates the RenderTemplateWithFuncs function
func sendTemplateWithFuncsMail(manager *mailstyler.Manager, receiver string) {
	log.Println("\nExample 3: Using RenderTemplateWithFuncs")

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

// sendTemplateWithCSSMail demonstrates the RenderTemplateWithCSS function
func sendTemplateWithCSSMail(manager *mailstyler.Manager, receiver string) {
	log.Println("\nExample 4: Using RenderTemplateWithCSS")
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

// sendTemplateWithFuncsAndCSSMail demonstrates the RenderTemplateWithFuncsAndCSS function
func sendTemplateWithFuncsAndCSSMail(manager *mailstyler.Manager, receiver string) {
	log.Println("\nExample 5: Using RenderTemplateWithFuncsAndCSS")
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

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize mailer manager
	manager := createManager(cfg)

	// Run examples
	sendTemplateMail(manager, cfg.MailReceiver)
	sendTemplateWithFuncsMail(manager, cfg.MailReceiver)
	sendTemplateWithCSSMail(manager, cfg.MailReceiver)
	sendTemplateWithFuncsAndCSSMail(manager, cfg.MailReceiver)
}
