package slack

import (
	"github.com/nlopes/slack"
)

func createSlackMessageOptions(alertName, alertText string, fields ...string) []slack.MsgOption {
	opts := make([]slack.MsgOption, 0)

	mainTextBlock := slack.NewTextBlockObject("mrkdwn", alertText, false, false)

	fieldsBlocks := make([]*slack.TextBlockObject, 0)

	for _, field := range fields {
		fieldsBlocks = append(fieldsBlocks, slack.NewTextBlockObject("mrkdwn", field, false, false))
	}

	sectionBlock := slack.NewSectionBlock(mainTextBlock, fieldsBlocks, nil)
	divBlock := slack.NewDividerBlock()
	alertNameBlock := slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", alertName, true, false))

	opts = append(opts, slack.MsgOptionBlocks(sectionBlock, divBlock, alertNameBlock))
	opts = append(opts, slack.MsgOptionBlocks(sectionBlock))
	opts = append(opts, slack.MsgOptionText(alertText, true), slack.MsgOptionAsUser(true))

	return opts
}
