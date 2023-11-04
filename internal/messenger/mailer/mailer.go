package mailer

import (
	"bytes"
	"crypto/tls"
	"html/template"

	"github.com/DMarinuks/disk-usage-warner/internal/logger"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/types"

	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

var (
	mailCfg *MailConfig
)

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

type Mailer struct {
	config MailConfig
	log    *zap.Logger
}

var _ types.Messenger = (*Mailer)(nil)

func New(cfg MailConfig) *Mailer {
	mailer := new(Mailer)
	mailer.config = cfg
	mailer.log = logger.Named("mailer")

	return mailer
}

func (m *Mailer) Send(hostname string, warnings []*types.WarningInfo) error {
	mail := gomail.NewMessage()

	mail.SetHeader("From", mailCfg.From)
	mail.SetHeader("To", mailCfg.Admins...)
	mail.SetHeader("Subject", mailCfg.Subject)

	body, err := loadTemplate(hostname, warnings)
	if err != nil {
		m.log.Error("error loading html template", zap.Error(err))
		return err
	}
	mail.SetBody("text/html; charset=UTF-8", body)

	d := gomail.NewDialer(mailCfg.Host, mailCfg.Port, mailCfg.Username, mailCfg.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{
		ServerName:         mailCfg.Host,
		InsecureSkipVerify: mailCfg.Insecure,
	}

	if err := d.DialAndSend(mail); err != nil {
		m.log.Error("error dial and send", zap.Error(err))
		return err
	}

	return nil
}

func loadTemplate(hostname string, warnings []*types.WarningInfo) (string, error) {
	templateData := struct {
		Host     string
		Warnings []*types.WarningInfo
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
