package main

import (
	"fmt"
	"os"

	"github.com/DMarinuks/disk-usage-warner/internal/diskchecker"
	"github.com/DMarinuks/disk-usage-warner/internal/logger"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/mailer"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

var cli struct {
	LogLevel string `help:"Log level" default:"warning" env:"DUW_LOG_LEVEL" short:"l"`
	// RocketChat config

	RocketChatToken  string `env:"ROCKET_CHAT_TOKEN"`
	RocketChatUserID string `env:"ROCKET_CHAT_USER_ID"`

	// Mailer config

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
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, "disk usage warner err: %v\n", err)
		os.Exit(1)
	}
}

func doMain() error {
	godotenv.Load()

	kongCtx := kong.Parse(&cli)

	if err := logger.Configure(cli.LogLevel); err != nil {
		return fmt.Errorf("error configuring logger: %w", err)
	}

	switch {
	case cli.RocketChatToken != "":
		diskchecker.SetDefaultMessenger(messenger.NewRocketChatMessenger(cli.RocketChatToken, cli.RocketChatUserID))
	case cli.MailPassword != "":
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

		diskchecker.SetDefaultMessenger(messenger.NewMailMessenger(mailerCfg))
	}

	if err := kongCtx.Run(); err != nil {
		return err
	}

	return nil
}
