package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailWarning(email, oldIP, newIP string) error {
	auth := smtp.PlainAuth("", email, os.Getenv("APP_PASSWORD"), "smtp.gmail.com")

	msg := fmt.Sprintf("Subject: App Warning: Your IP Address has been changed\r\n"+"Old IP: %s\r\n"+"New IP: %s\r\n"+"If it wasn't you, please secure your account\r\n", oldIP, newIP)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "danil.novo86@gmail.com", []string{email}, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
