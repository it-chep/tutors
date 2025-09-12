package smtp

import (
	"gopkg.in/gomail.v2"
)

type ClientSmtp struct {
	dialer  *gomail.Dialer
	address string
}

func NewClientSmtp(address, smtpSecret string) *ClientSmtp {
	d := gomail.NewDialer("smtp.gmail.com", 587, address, smtpSecret)
	return &ClientSmtp{
		dialer:  d,
		address: address,
	}
}

type EmailParams struct {
	Body        string
	Destination string
	Subject     string
}

func (c *ClientSmtp) SendEmail(params EmailParams) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.address)
	msg.SetHeader("To", params.Destination)
	msg.SetHeader("Subject", params.Subject)
	msg.SetBody("text/html", params.Body)

	if err := c.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
