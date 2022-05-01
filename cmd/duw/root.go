package main

import (
	"DMarinuks/disk-usage-warner/logger"
	"DMarinuks/disk-usage-warner/mailer"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

var cli struct {
	LogLevel string `help:"Log level" default:"warning" env:"DUW_LOG_LEVEL" short:"l"`

	Admins       []string `help:"Email addresses of the admins" env:"DUW_ADMINS"`
	MailSubject  string   `default:"Disk Usage Warning" env:"DUW_MAIL_SUBJECT"`
	MailFrom     string   `env:"DUW_MAIL_FROM"`
	MailHost     string   `env:"DUW_MAIL_HOST"`
	MailPort     int      `env:"DUW_MAIL_PORT"`
	MailUsername string   `env:"DUW_MAIL_USERNAME"`
	MailPassword string   `env:"DUW_MAIL_PASSWORD"`
	MailTLSSkip  bool     `help:"Skip tls verify" env:"DUW_MAIL_TLS_SKIP"`

	Check runCheck `cmd:""`
}

func main() {
	if err := DoMain(); err != nil {
		fmt.Fprintf(os.Stderr, "disk usage warner err: %v\n", err)
		os.Exit(1)
	}
}

func DoMain() error {
	// load env file
	godotenv.Load()

	// parse cli
	kongCtx := kong.Parse(&cli)

	// configure logging
	if err := logger.Configure(cli.LogLevel); err != nil {
		return fmt.Errorf("error configuring logger: %w", err)
	}

	// configure mailer
	mailerCfg := mailer.MailConfig{
		From:     cli.MailFrom,
		Subject:  cli.MailSubject,
		Admins:   cli.Admins,
		Host:     cli.MailHost,
		Port:     cli.MailPort,
		Username: cli.MailUsername,
		Password: cli.MailPassword,
		Insecure: cli.MailTLSSkip,
	}
	mailer.Configure(mailerCfg)

	// run command
	if err := kongCtx.Run(); err != nil {
		return err
	}

	return nil
}
