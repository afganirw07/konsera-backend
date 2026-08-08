package utils

import (
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"os"
	"github.com/joho/godotenv"
)

type EmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

var (
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
)