package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
)

type EmailData struct {
	TotalBalance string
	StatsByMonth map[string]MonthlyStats
}

func sendEmail(stats TransactionStats, recipient string) error {
	emailBody, err := generateHTMLEmail(stats)
	if err != nil {
		return fmt.Errorf("error creating the email template: %v", err)
	}

	gmailUsername := os.Getenv("GMAIL_USERNAME")
	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	gmailAuth := smtp.PlainAuth("", gmailUsername, gmailPassword, "smtp.gmail.com")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	defer writer.Close()

	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("error creating HTML part: %v", err)
	}
	htmlPart.Write([]byte(emailBody))

	logoPath := filepath.Join("templates", "stori-logo.png")
	logoData, err := os.ReadFile(logoPath)
	if err == nil {
		imagePart, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {"image/png"},
			"Content-Transfer-Encoding": {"base64"},
			"Content-ID":                {"<stori-logo>"},
			"Content-Disposition":       {"inline"},
		})
		if err != nil {
			return fmt.Errorf("error creating image part: %v", err)
		}

		encoder := base64.NewEncoder(base64.StdEncoding, imagePart)
		encoder.Write(logoData)
		encoder.Close()
	}

	var emailBuffer bytes.Buffer

	emailBuffer.WriteString(fmt.Sprintf("From: %s\r\n", gmailUsername))
	emailBuffer.WriteString(fmt.Sprintf("To: %s\r\n", recipient))
	emailBuffer.WriteString("Subject: Your Stori Account Transaction Summary\r\n")
	emailBuffer.WriteString("MIME-Version: 1.0\r\n")
	emailBuffer.WriteString(fmt.Sprintf(
		"Content-Type: multipart/related; boundary=%s\r\n\r\n", writer.Boundary()))
	emailBuffer.WriteString("\r\n")

	emailBuffer.Write(body.Bytes())

	mailErr := smtp.SendMail(
		"smtp.gmail.com:587", gmailAuth, gmailUsername, []string{recipient}, emailBuffer.Bytes())
	if mailErr != nil {
		return fmt.Errorf("error sending email: %v", mailErr)
	}
	return nil
}

func generateHTMLEmail(stats TransactionStats) (string, error) {
	funcMap := template.FuncMap{
		"monthName":      monthName,
		"add":            func(a, b int) int { return a + b },
		"formatCurrency": formatCurrency,
		"gt":             func(a, b int) bool { return a > b },
	}

	templatePath := filepath.Join("templates", "email-template.html")
	htmlTemplate, err := template.New("email").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error loading template file: %v", err)
	}

	data := EmailData{
		TotalBalance: formatCurrency(stats.TotalBalance),
		StatsByMonth: stats.StatsByMonth,
	}

	var buf bytes.Buffer
	if err := htmlTemplate.ExecuteTemplate(&buf, "email-template.html", data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}

func monthName(month string) string {
	months := map[string]string{
		"01": "January", "02": "February", "03": "March", "04": "April",
		"05": "May", "06": "June", "07": "July", "08": "August",
		"09": "September", "10": "October", "11": "November", "12": "December",
	}

	if name, ok := months[month]; ok {
		return name
	}
	return month
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
