package messenger

import (
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/mailer"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/rocketchat"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/types"
)

func NewRocketChatMessenger(token, userID string) types.Messenger {
	return rocketchat.New(token, userID)
}

func NewMailMessenger(cfg mailer.MailConfig) types.Messenger {
	return mailer.New(cfg)
}
