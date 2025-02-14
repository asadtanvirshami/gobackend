package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// Send OTP Email
func SendOTPEmail(email, otp string) {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))
	to := []string{email}
	msg := []byte("Subject: Email Verification\n\nYour OTP is: " + otp)

	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("SMTP_USER"), to, msg)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
}
