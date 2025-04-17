// Package mail provides a simple and flexible way to send emails using SMTP,
// with support for rendering HTML templates and dynamic data.
//
// Example usage:
//
//	manager, err := mail.NewManager("smtp.example.com", "587", "you@example.com", "./templates")
//	if err != nil {
//	    log.Fatalf("failed to create manager: %v", err)
//	}
//
//	body, err := manager.RenderTemplate("hello.html", map[string]any{
//	    "Name": "Someone",
//	})
//	if err != nil {
//	    log.Fatalf("failed to render template: %v", err)
//	}
//
//	err = manager.SendMail(mail.MailMessage{
//	    Subject: "Hello",
//	    Message: body,
//	    To:      []string{"someone@example.com"},
//	})
//	if err != nil {
//	    log.Fatalf("failed to send mail: %v", err)
//	}
package mail
