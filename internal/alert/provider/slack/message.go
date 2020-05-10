package slack

import (
	"github.com/nlopes/slack"
)

func createSlackMessageOptions(alertText, imageURL string, fields ...string) []slack.MsgOption {
	blocks := make([]slack.Block, 0)

	if imageURL != "" {
		imgTextBlock := slack.NewTextBlockObject("plain_text", "chart", false, false)
		imgBlock := slack.NewImageBlock(imageURL, "chart", "", imgTextBlock)
		blocks = append(blocks, imgBlock)
	}

	if alertText != "" {
		mainTextBlock := slack.NewTextBlockObject("mrkdwn", alertText, false, false)
		mainSectionBlock := slack.NewSectionBlock(mainTextBlock, nil, nil)
		blocks = append(blocks, mainSectionBlock)
	}

	fieldsBlocks := make([][]*slack.TextBlockObject, 0)

	ff := make([]*slack.TextBlockObject, 0)

	// Slack supports up to 10 fields
	for idx, field := range fields {
		ff = append(ff, slack.NewTextBlockObject("mrkdwn", field, false, false))

		if (idx+1)%10 == 0 {
			fieldsBlocks = append(fieldsBlocks, ff)
			ff = make([]*slack.TextBlockObject, 0)
		}
	}

	if len(fields) > 0 && len(fields)%10 != 0 {
		fieldsBlocks = append(fieldsBlocks, ff)
	}

	for _, fieldBlock := range fieldsBlocks {
		blocks = append(blocks, slack.NewSectionBlock(nil, fieldBlock, nil))
	}

	opts := make([]slack.MsgOption, 0)
	opts = append(opts, slack.MsgOptionBlocks(blocks...))

	if alertText != "" {
		opts = append(opts, slack.MsgOptionAsUser(true))
	}

	return opts
}
