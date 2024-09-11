package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail sends an email using the SMTP configuration from environment variables.
func GoogleSendEmail(to string, subject string, body string, link string) error {
	from := os.Getenv("MAIL_FROM_ADDRESS")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	username := os.Getenv("MAIL_USERNAME")

	fmt.Println("MAIL_FROM_ADDRESS:", from)
	fmt.Println("MAIL_USERNAME:", username)
	fmt.Println("MAIL_PASSWORD:", password)
	fmt.Println("MAIL_HOST:", smtpHost)
	fmt.Println("MAIL_PORT:", smtpPort)

	// Convert smtpPort from string to int
	port := 2525 // Default Mailtrap port for example, replace if necessary
	fmt.Sscanf(smtpPort, "%d", &port)

	// Use fmt.Sprintf to build clean HTML body without showing the raw link
	htmlBody := fmt.Sprintf(`
<html>
<body>
	<p>Please verify your email by clicking the button below:</p>
	<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-align: center; text-decoration: none; border-radius: 5px;">
		Verify Email
	</a>
</body>
</html>
`, link)

	// Create a new gomail message
	m := gomail.NewMessage()

	// Set the sender address
	m.SetHeader("From", from)

	// Set the recipient address
	m.SetHeader("To", to)

	// Set the subject of the email
	m.SetHeader("Subject", subject)

	// Set the body of the email to HTML format
	m.SetBody("text/html", htmlBody)

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, port, username, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to:", to)
	return nil
}
