package service

import (
	"advanced-webapp-project/config"
	"advanced-webapp-project/model"
	"advanced-webapp-project/utils"
	"bytes"
	"crypto/tls"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
)

type IMailerService interface {
	SendMail(user *model.User, msg *Message)
}

type mail struct {
	from     string
	to       string
	smtpUser string
	smtpPass string
	smtpHost string
	smtpPort int
	logger   *utils.Logger
}

type Message struct {
	URL      string
	FullName string
	Subject  string
}

func NewMailerService(logger *utils.Logger) *mail {
	return &mail{
		logger: logger,
	}
}

func (m *mail) parseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func (m *mail) SendMail(user *model.User, msg *Message) {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		m.logger.Error(err.Error())
		return
	}

	m.from = "admin@kameyoko.com"
	m.to = user.Email
	m.smtpUser = appConfig.SMTPUser
	m.smtpPass = appConfig.SMTPPass
	m.smtpHost = appConfig.SMTPHost
	m.smtpPort = appConfig.SMTPPort

	var body bytes.Buffer
	tpl, err := m.parseTemplateDir("templates")
	if err != nil {
		m.logger.Error(err.Error())
		return
	}

	tpl.ExecuteTemplate(&body, "verificationCode.gohtml", &msg)

	email := gomail.NewMessage()
	email.SetHeader("From", m.from)
	email.SetHeader("To", m.to)
	email.SetHeader("Subject", msg.Subject)
	email.SetBody("text/html", body.String())
	email.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dial := gomail.NewDialer(appConfig.SMTPHost, appConfig.SMTPPort, appConfig.SMTPUser, appConfig.SMTPPass)
	dial.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = dial.DialAndSend(email); err != nil {
		m.logger.Error(err.Error())
		return
	}
}
