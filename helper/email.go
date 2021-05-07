package helper

import (
	"log"
	"net/smtp"
)

func SendEmail(body string) {
	from := ""
	pass := ""
	to := []string{""}

	msg := "From: " + from + "\n" +
		"Subject: Vaccine Notification\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		   from, to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}