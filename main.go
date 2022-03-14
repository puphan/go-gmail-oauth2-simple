package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type DataToHTML struct {
	Name    string
	Address string
}

func main() {
	dataToHTML := new(DataToHTML)
	dataToHTML.Name = "Hello World"
	dataToHTML.Address = "555/777"

	OAuthGmailService()
	sendStatus, err := SendEmailOAUTH2("to_email", dataToHTML, "testTemplate.html")
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	log.Println("status send mail:", sendStatus)

}

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  "your_access_token",
		RefreshToken: "your_refresh_token",
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
		panic(err)
	}

	GmailService = srv
	if GmailService != nil {
		log.Println("Email service is initialized")
	}

}

func SendEmailOAUTH2(to string, data interface{}, template string) (bool, error) {

	emailBody, err := parseTemplate(template, data)
	if err != nil {
		log.Println(err.Error())
		return false, errors.New("unable to parse email template")
	}

	//or text html
	// htmlText := `<html><header></header><body><h1>Test 01</h1><h2>eiei</h2><h2>123/22</h2></body></html>`
	// emailBody := htmlText

	var message gmail.Message

	email := "your_email_sender"
	displayName := "ECO SYSTEM"
	fromMail := fmt.Sprintf("From: %s <%s> \r\n", displayName, email)
	emailTo := "To: " + to + "\r\n"
	subjectText := "Test 123 Email form Gmail API using OAuth2 (utf-8 ทดสอบ)"
	subject := fmt.Sprintf("Subject: =?utf-8?B?%s?=\n", base64.StdEncoding.EncodeToString([]byte(subjectText)))
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fromMail + emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)
	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("emailTemplates/%s", templateFileName))
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Println(err.Error())
		return "", err
	}
	body := buf.String()
	return body, nil
}
