package slack_integration

import (
	"github.com/bluele/slack"
	"os"
)

func SendSlack(message string) {
	slackApi := slack.New(os.Getenv("SLACK_TOKEN"))
	slackApi.ChatPostMessage("exchange-notification", message, nil)
}
