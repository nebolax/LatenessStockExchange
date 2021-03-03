package main

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func send() {
	from := mail.NewEmail("Example User", "lateness.se.mailer@gmail.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "s22e_nebolsin@179.ru")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>Hello pesiy pes</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println("error")
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
