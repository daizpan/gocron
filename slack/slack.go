package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackClient struct {
	WebHookUrl string
}

type SlackMessage struct {
	Text string `json:"text,omitempty"`
}

func NewSlackClient(webHookUrl string) *SlackClient {
	return &SlackClient{
		WebHookUrl: webHookUrl,
	}
}

func (sc *SlackClient) Post(msg string) (err error) {
	slackMsg := SlackMessage{
		Text: msg,
	}
	msgBody, _ := json.Marshal(slackMsg)
	res, err := http.Post(sc.WebHookUrl, "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		return err
	}
	defer func() { err = res.Body.Close() }()

	return nil
}
