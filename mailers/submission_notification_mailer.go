package mailers

import (
	"github.com/flosch/pongo2"
	"github.com/plentiform/plentiform/models"
	"github.com/sendgrid/rest"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	SUBMISSION_NOTIFICATION_TEXT_TEMPLATE = `TODO Plaintext Emails`

	SUBMISSION_NOTIFICATION_HTML_TEMPLATE = `
<h1>New Form Submission</h1> 
<br>
<pre>{{ json }}</pre>`
)

func SendSubmissionNotification(emailClient *sendgrid.Client, user *models.User, form *models.Form, json []byte) (*rest.Response, error) {
	c := pongo2.Context{"json": string(json)}

	htmlContent, _ := pongo2.FromString(SUBMISSION_NOTIFICATION_HTML_TEMPLATE)
	html, _ := htmlContent.Execute(c)

	plainTextContent, _ := pongo2.FromString(SUBMISSION_NOTIFICATION_TEXT_TEMPLATE)
	plainText, _ := plainTextContent.Execute(c)

	message := mail.NewSingleEmail(
		mail.NewEmail("Plentiform", "team@plentiform.com"),
		"New Form Submission",
		mail.NewEmail(user.Name, user.Email),
		plainText,
		html,
	)

	return emailClient.Send(message)
}
