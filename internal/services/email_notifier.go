package services

import (
	"time"

	"github.com/spf13/viper"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailNotifier struct {
	Emails  []string
	Config  *viper.Viper
	Subject string
}

func NewEmailNotifier(config *viper.Viper, emails []string, subject string) *EmailNotifier {
	return &EmailNotifier{
		Config:  config,
		Emails:  emails,
		Subject: subject,
	}
}

func (e *EmailNotifier) SendNotification(message string) error {

	server := mail.NewSMTPClient()
	server.Host = e.Config.GetString("SMTP_HOST")
	server.Port = e.Config.GetInt("SMTP_PORT")
	server.Username = e.Config.GetString("SMTP_USERNAME")
	server.Password = e.Config.GetString("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpCliente, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("Roberto Suárez <electrosonix12@gmail.com>").
		AddTo(e.Emails...).
		SetSubject(e.Subject).
		SetListUnsubscribe("<mailto:unsubscribe@example.com?subject=https://example.com/unsubscribe>")

	email.SetBody(mail.TextPlain, message)

	err = email.Send(smtpCliente)
	if err != nil {
		return err
	}

	return nil
}

func (e *EmailNotifier) ResetPassword(message string, url string) error {
	server := mail.NewSMTPClient()
	server.Host = e.Config.GetString("SMTP_HOST")
	server.Port = e.Config.GetInt("SMTP_PORT")
	server.Username = e.Config.GetString("SMTP_USERNAME")
	server.Password = e.Config.GetString("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpCliente, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("Roberto Suárez <electrosonix12@gmail.com>").
		AddTo(e.Emails...).
		SetSubject(e.Subject).
		SetListUnsubscribe("<mailto:unsubscribe@example.com?subject=https://example.com/unsubscribe>")

	email.SetBody(mail.TextHTML, message)

	err = email.Send(smtpCliente)
	if err != nil {
		return err
	}

	return nil
}
