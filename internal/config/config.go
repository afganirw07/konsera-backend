package config

import "os"

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromName     string
	FromEmail    string
}

func LoadEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromName:     os.Getenv("SMTP_FROM_NAME"),
		FromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
	}
}