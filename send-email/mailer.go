package main

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

// Mailer is a struct which takes a config and includes a Send method for sending mail
type Mailer struct {
	EmailConfig
}

// Send sends an email using the Mailer's config
func (m Mailer) Send() error {
	client, err := smtp.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return err
	}

	client.StartTLS(&tls.Config{ServerName: m.Host})

	if err = client.Auth(smtp.PlainAuth("", m.From, m.Password, m.Host)); err != nil {
		return err
	}

	if err = client.Mail(m.From); err != nil {
		return err
	}

	if err = client.Rcpt(m.To); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(m.makeMessage()))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}

func (m Mailer) makeMessage() (message string) {
	from := mail.Address{Name: "Lauren Settembrino", Address: m.From}
	to := mail.Address{Name: "Lauren Settembrino", Address: m.To}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = m.Subject
	headers["Reply-To"] = m.SenderEmail

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += fmt.Sprintf(
		"\r\nNew message received from %s via laurensettembrino.com:\n\n%s", m.SenderName, m.Message,
	)

	return message
}
