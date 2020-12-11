package slack

import (
	"os"
	"testing"
)

func TestSlackClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		webHookUrl string
		msg        string
	}{
		{
			name:       "Post-test",
			webHookUrl: os.Getenv("SLACK_URL"),
			msg:        "Post test!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewSlackClient(tt.webHookUrl)
			if err := sc.Post(tt.msg); err != nil {
				t.Errorf("SlackClient.Post() error = %v", err)
			}
		})
	}
}
