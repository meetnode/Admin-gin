package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

var url = os.Getenv("URL")

func SendMail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	addr := host + ":" + port
	auth := smtp.PlainAuth("", from, pass, host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

func SendVerificationEmail(to, token string) error {
	verifyLink := fmt.Sprintf("%s/api/verify?token=%s", url, token)
	subject := "Verify your email"
	body := "Click here to verify your email: " + verifyLink
	return SendMail(to, subject, body)
}

func SendPasswordResetEmail(to, token string) error {
	resetLink := fmt.Sprintf("%s/api/reset-password?token=%s", url, token)
	subject := "Reset your password"
	body := "Click here to reset your password: " + resetLink
	return SendMail(to, subject, body)
}
