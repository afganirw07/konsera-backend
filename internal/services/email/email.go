package email

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"log"
	"net/smtp"
	"time"

	dto "konsera-backend/internal/DTO/email"
	templates "konsera-backend/internal/templates"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromName     string
	FromEmail    string
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) buildEmailFromTemplate(templateName string, data interface{}) (string, error) {
	templatePath := fmt.Sprintf("gmail/%s", templateName)

	funcMap := htmltemplate.FuncMap{
		"formatRupiah": func(amount float64) string {
			p := message.NewPrinter(language.Indonesian)
			return p.Sprintf("Rp%.0f", amount)
		},
	}

	tmpl, err := htmltemplate.
		New(templateName).
		Funcs(funcMap).
		ParseFS(templates.EmailTemplates, templatePath)

	if err != nil {
		log.Printf(
			"[EMAIL_SERVICE] Failed to parse template %s: %v",
			templateName,
			err,
		)

		return "", fmt.Errorf(
			"failed to parse email template: %w",
			err,
		)
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, data); err != nil {
		log.Printf(
			"[EMAIL_SERVICE] Failed to execute template %s: %v",
			templateName,
			err,
		)

		return "", fmt.Errorf(
			"failed to execute email template: %w",
			err,
		)
	}

	return output.String(), nil
}

func (s *Service) SendEmail(to string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUsername,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	msg := []byte(
		"From: " + s.config.FromName + " <" + s.config.FromEmail + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			htmlBody,
	)

	address := fmt.Sprintf(
		"%s:%s",
		s.config.SMTPHost,
		s.config.SMTPPort,
	)

	start := time.Now()

	log.Printf("[EMAIL_SERVICE] Sending email to %s...", to)

	err := smtp.SendMail(
		address,
		auth,
		s.config.FromEmail,
		[]string{to},
		msg,
	)

	log.Printf(
		"[EMAIL_SERVICE] SMTP finished in %v",
		time.Since(start),
	)

	if err != nil {
		log.Printf(
			"[EMAIL_SERVICE] Failed to send email to %s: %v",
			to,
			err,
		)

		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf(
		"[EMAIL_SERVICE] Email successfully sent to %s",
		to,
	)

	return nil
}

func (s *Service) SendRegisterOTP(to string, data dto.RegisterOTPEmailData) error {
	htmlBody, err := s.buildEmailFromTemplate("otp_register.html", data)
	if err != nil {
		return err
	}

	subject := "Verify your email - Konsera"

	if err := s.SendEmail(
		to,
		subject,
		htmlBody,
	); err != nil {
		return err
	}

	return nil
}
