package helper

import (
	"log"
	"net/smtp"
)

func SendEmail(body string) {
	from := "pratik.narode@velotio.com"
	pass := "wqnenqzuiachflnn"
	to := []string{"ankitadaher@gmail.com", "pratiknarode143@gmail.com", "shashanknarode133@gmail.com"}

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