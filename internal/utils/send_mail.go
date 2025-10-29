package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var url = os.Getenv("URL")

func SendMail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	addr := host + ":" + port
	auth := smtp.PlainAuth("", from, pass, host)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

func SendVerificationEmail(to, token string) error {
	verifyLink := fmt.Sprintf("%s/api/verify?token=%s", url, token)
	subject := "Verify your email"
	body := "Click here to verify your email: " + verifyLink
	return SendMail(to, subject, body)
}

func SendPasswordResetEmail(to, token string) error {
	baseURL := os.Getenv("URL")
	resetLink := fmt.Sprintf("%s/api/reset-password?token=%s", baseURL, token)

	subject := "Reset your password"

	templatePath := "internal/templates/reset_password.html"
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading email template: %w", err)
	}
	body := string(content)

	//Add more replacements if needed in future
	replacements := map[string]string{
		"RESET_LINK": resetLink,
	}
	for key, value := range replacements {
		placeholder := fmt.Sprintf("%%%s%%", key)
		body = strings.ReplaceAll(body, placeholder, value)
	}

	return SendMail(to, subject, body)
}
