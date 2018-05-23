package email

import (
	"bytes"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/integration/sendgrid_service"
	"html/template"
	"os"
)

func SendSystemEmailWithTemplate(toName string, toAddress string, language string, subject string, templateKey string, data interface{}) error {
	fromName := os.Getenv("EMAIL_FROM_NAME")
	fromAddress := os.Getenv("EMAIL_FROM_ADDRESS")

	err := SendEmailWithTemplate(fromName, fromAddress, toName, toAddress, language, subject, templateKey, data)
	return err
}

func SendEmailWithTemplate(fromName string, fromAddress string, toName string, toAddress string, language string, subject string, templateKey string, data interface{}) error {
	templatePath := "./templates/" + TemplateName[templateKey] + language + ".html"

	t, _ := template.ParseFiles(templatePath)
	buffer := bytes.NewBufferString("")
	err := t.Execute(buffer, data)

	if err == nil {
		err = sendgrid_service.SendEmail(fromName, fromAddress, toName, toAddress, subject, buffer.String())
	} else {
		err = api_error.PropagateError(api_error.ReadTemplateError, err)
	}

	return err
}
