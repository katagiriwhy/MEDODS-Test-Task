package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailWarning(email, oldIP, newIP string) error {
	const from = "danil.novo86@gmail.com"

	auth := smtp.PlainAuth("", email, os.Getenv("APP_PASSWORD"), "smtp.gmail.com")

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: App Warning: Your IP Address has been changed\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"Security Alert!\r\n\r\n"+
			"Your account was accessed from a new IP address:\r\n"+
			"Old IP: %s\r\n"+
			"New IP: %s\r\n\r\n"+
			"If this wasn't you, please secure your account!\r\n",
		from, email, oldIP, newIP)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
