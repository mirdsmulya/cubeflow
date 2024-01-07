package slack

import (
	"fmt"

	"cubeflow/pkg/config"

	"github.com/slack-go/slack"
)

func SendSlackMessage(slackChannelID string, message string) {
	slackToken := config.Variable.Slack.Token

	api := slack.New(slackToken)
	_, _, err := api.PostMessage(slackChannelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("Error sending message to Slack: %v\n", err)
		return
	}
}