package app_email

import (
	"net/smtp"

	"github.com/tidinio/src/component/configuration"
	"github.com/tidinio/src/component/template"
)

const (
	EmailInfo = "newpsel@gmail.com"
)

type Request struct {
	body    string
	from    string
	subject string
	to      []string
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    "",
	}
}

func (r *Request) ParseEmailTemplate(templateFileName string, templateData interface{}) error {
	body, err := app_template.ParseEmailTemplate(templateFileName, templateData)
	r.body = body
	return err
}

func (r *Request) SetFrom(emailFrom string) {
	r.from = emailFrom
}

func getAuth() smtp.Auth {
	mailerUser, _ := app_conf.Data.String("smtp.mailerUser")
	mailerPassword, _ := app_conf.Data.String("smtp.mailerPassword")
	return smtp.PlainAuth("", mailerUser, mailerPassword, "smtp.gmail.com")
}

func (r *Request) SendEmail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"
	if err := smtp.SendMail(addr, getAuth(), r.from, r.to, msg); err != nil {
		return err
	}
	return nil
}
