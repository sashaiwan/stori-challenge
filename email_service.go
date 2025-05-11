package main

import (
	"bytes"
	"fmt"
	"html/template"
)

type EmailData struct {
	TotalBalance string
	StatsByMonth map[string]MonthlyStats
}

// TODO:
// func sendEmail()

func generateHTMLEmail(stats TransactionStats) (string, error) {
	funcMap := template.FuncMap{
		"monthName":      monthName,
		"add":            func(a, b int) int { return a + b },
		"formatCurrency": formatCurrency,
		"gt":             func(a, b int) bool { return a > b },
	}

	template, err := template.New("email").Funcs(funcMap).ParseFiles("./templates/email-template.html")
	if err != nil {
		return "", fmt.Errorf("error loading template file: %v", err)
	}

	data := EmailData{
		TotalBalance: formatCurrency(stats.TotalBalance),
		StatsByMonth: stats.StatsByMonth,
	}

	var buf bytes.Buffer
	if err := template.ExecuteTemplate(&buf, "email-template.html", data); err != nil {
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
