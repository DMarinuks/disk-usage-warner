package rocketchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DMarinuks/disk-usage-warner/internal/messenger/types"
)

var _ types.Messenger = (*RocketChatMessenger)(nil)

type RocketChatMessenger struct {
	token  string
	userID string
}

func New(token, userID string) *RocketChatMessenger {
	messenger := new(RocketChatMessenger)
	messenger.token = token
	messenger.userID = userID

	return messenger
}

func (m *RocketChatMessenger) Send(hostname string, warnings []*types.WarningInfo) error {
	var msg strings.Builder

	msg.WriteString("#### Disk Usage Warning\n")
	msg.WriteString("Host: `" + hostname + "`\n")

	for _, warning := range warnings {
		msg.WriteString("Mount: `" + warning.Device + "` used " + warning.Percent + "\n")
	}

	return m.sendRocketChatMessage(msg.String())
}

func (m *RocketChatMessenger) sendRocketChatMessage(msg string) error {
	rcMsg := struct {
		Alias   string `json:"alias"`
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}{
		Channel: "#gfbio-monitoring",
		Text:    msg,
	}

	payload, err := json.Marshal(rcMsg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://chat.gwdg.de/api/v1/chat.postMessage", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", m.token)
	req.Header.Add("X-User-Id", m.userID)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("body", string(body))

		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
