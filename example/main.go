package main

import (
	"log"

	"github.com/phzeng0726/go-mailer/mail"
)

func main() {
	manager, err := mail.NewManager("smtp.example.com", "587", "you@example.com", "./templates")
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}

	body, err := manager.RenderTemplate("hello.html", map[string]any{
		"Name": "Someone",
	})
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}

	err = manager.SendMail(mail.MailMessage{
		Subject: "Hello",
		Message: body,
		To:      []string{"someone@example.com"},
	})
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}

	log.Println("Mail sent successfully!")
}
