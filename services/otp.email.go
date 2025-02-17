package services

import (
	gomail "gopkg.in/gomail.v2"
)

// Send OTP Email
func SendOTPEmail(email, otp string) {
	from := "asadtanvir20@gmail.com"
	to := email
	host := "smtp-relay.brevo.com"
	port := 587

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", " Email Verification OTP")
	// text/html for a html email
	msg.SetBody("text/html", "<div><p>Do not share this passcode with anyone else. <br/></p><strong>"+otp+"</strong></div>")

	n := gomail.NewDialer(host, port, from, "ZGRIB7TKYg083z9h")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}
