package service

import (
	"f_blog/backend/internal/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(cfg *config.EmailConfig, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}
