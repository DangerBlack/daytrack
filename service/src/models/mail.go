package models

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"text/template"

	"512b.it/daytrack/src/utils"
	gomail "gopkg.in/mail.v2"
)

//go:embed magic_link.html
var magicLinkTemplate string

type MagicLinkMail struct {
	Url string
}

func SendMagicLink(ctx context.Context, configuration Configuration, to string, magicLink string) error {
	var err error
	i := MagicLinkMail{
		Url: magicLink,
	}
	logger := utils.InitServiceLogger("Mail")
	t := template.New("")

	t, err = t.Parse(magicLinkTemplate)
	if err != nil {
		logger(ctx).Err(err).Msg("Failed to parse template")
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		logger(ctx).Err(err).Msg("Failed to execute template")
	}
	result := tpl.String()

	return SendMail(configuration, to, "Magic Link", "This is a magic link: "+magicLink, result)
}

func SendMail(configuration Configuration, to string, subject string, bodyText string, bodyHTML string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", configuration.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", bodyText)
	m.SetBody("text/html", bodyHTML)

	// Settings for SMTP server
	d := gomail.NewDialer(configuration.Mail.SMTP, configuration.Mail.PORT, configuration.Mail.From, configuration.Mail.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	return nil
}
