package main

import (
	"crypto/tls"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

func sendEmail(title string, body string, toEmail string) {

	from := os.Getenv("EmailerEmail")
	password := os.Getenv("EmailerEmailPassward")

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", toEmail)

	// Set E-Mail subject
	m.SetHeader("Subject", title)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
