package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/DMarinuks/disk-usage-warner/logger"

	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

var (
	log     *zap.Logger
	mailCfg *MailConfig
)

type WarningInfo struct {
	Device  string
	Percent string
}

type MailConfig struct {
	From     string
	Subject  string
	Admins   []string
	Host     string
	Port     int
	Username string
	Password string
	Insecure bool
}

func Configure(cfg MailConfig) {
	log = logger.Named("mailer")
	mailCfg = &cfg
}

func SendMail(warnings []*WarningInfo) error {
	hostname, _ := os.Hostname()
	hostname = strings.ToLower(strings.TrimSpace(hostname))
	if len(hostname) == 0 {
		return fmt.Errorf("empty hostname is invalid")
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", mailCfg.From)

	// Set E-Mail receivers
	m.SetHeader("To", mailCfg.Admins...)

	// Set E-Mail subject
	m.SetHeader("Subject", mailCfg.Subject)

	// Set E-Mail body.
	body, err := loadTemplate(hostname, warnings)
	if err != nil {
		log.Error("error loading html template", zap.Error(err))
		return err
	}
	m.SetBody("text/html; charset=UTF-8", body)

	// Settings for SMTP server
	d := gomail.NewDialer(mailCfg.Host, mailCfg.Port, mailCfg.Username, mailCfg.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{
		ServerName:         mailCfg.Host,
		InsecureSkipVerify: mailCfg.Insecure,
	}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Error("error dial and send", zap.Error(err))
		return err
	}

	return nil
}

func loadTemplate(hostname string, warnings []*WarningInfo) (string, error) {
	templateData := struct {
		Host     string
		Warnings []*WarningInfo
	}{
		Host:     hostname,
		Warnings: warnings,
	}
	t, err := template.New("warning").Parse(warningHTML)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}
