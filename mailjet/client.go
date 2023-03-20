package mailjet

import (
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

func SendVerification(email string, code string) error {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "no-reply@do-notes.pp.ua",
				Name:  "do-notes",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  "user",
				},
			},
			Subject:  "do-notes - verify your email",
			TextPart: "Email verification code is " + code,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)

	return err
}
