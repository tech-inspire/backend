package mail

import (
	"fmt"
	"net/smtp"

	"github.com/tech-inspire/backend/auth-service/internal/config"
)

type Client struct {
	auth smtp.Auth
	from string
	addr string
}

func NewClient(cfg *config.Config) (*Client, error) {
	// Connect to the SMTP server.
	auth := smtp.PlainAuth("", cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host)

	addr := fmt.Sprintf("%s:%s", cfg.SMTP.Host, cfg.SMTP.Port)

	return &Client{
		auth: auth,
		from: cfg.SMTP.From,
		addr: addr,
	}, nil
}

const mailTemplate = "To: %s\r\n" +
	"Subject: %s\r\n" +
	"\r\n" +
	"%s\r\n"

func (c *Client) SendMail(to string, message Message) error {
	msg := fmt.Sprintf(mailTemplate, to, message.Title, message.Message)

	err := smtp.SendMail(c.addr, c.auth, c.from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("smtp.SendMail: %w", err)
	}

	return nil
}
