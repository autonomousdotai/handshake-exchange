package sendgrid_service

import (
	"github.com/duyhtq/crypto-exchange-service/api_error"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

func SendEmail(fromName string, fromAddress string, toName string, toAddress string, subject string, body string) error {
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	from := mail.NewEmail(fromName, fromAddress)
	to := mail.NewEmail(toName, toAddress)

	message := mail.NewSingleEmail(from, subject, to, body, body)

	_, err := client.Send(message)
	err = api_error.PropagateError(api_error.SendGridError, err)

	return err
}
