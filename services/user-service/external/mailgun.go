package external

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/moshrank/spacey-backend/config"
)

type EmailSenderInterface interface {
	SendEmail(recipient, validationLink string) error
}

type EmailSender struct {
	mg  mailgun.Mailgun
	cfg config.ConfigInterface
}

func NewEmailSender(cfg config.ConfigInterface) EmailSenderInterface {
	domain := cfg.GetDomain()
	mg := mailgun.NewMailgun(domain, cfg.GetMailGunAPIKey())
	mg.SetAPIBase(mailgun.APIBaseEU)

	return &EmailSender{mg: mg, cfg: cfg}
}

func (e *EmailSender) SendEmail(recipient, validationLink string) error {

	domain := e.cfg.GetDomain()
	sender := fmt.Sprintf("noreply@<%s>", domain)
	subject := "Welcome to Spacey! Please Validate your email."
	body := ""

	message := e.mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("email_validation")
	err := message.AddTemplateVariable(
		"email_validation_link",
		validationLink,
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := e.mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}
