package main

import (
	"log"
	"net/smtp"
)

func main() {
	addr := "localhost:1025"
	from := "sender@example.com"
	to := []string{"alice@example.net"}
	msg := []byte("From: Sender <sender@example.com>\r\n" +
		"To: Alice <alice@example.net>\r\n" +
		"Subject: Hello from Go test SMTP\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		"This is a test email body.\r\n")

	// No auth needed for our test server
	if err := smtp.SendMail(addr, nil, from, to, msg); err != nil {
		log.Fatal(err)
	}
}
