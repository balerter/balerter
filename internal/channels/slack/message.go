package slack

import (
	"github.com/nlopes/slack"
)

func createSlackMessageOptions(alertText, imageURL string) []slack.MsgOption {
	opts := make([]slack.MsgOption, 0)
	blocks := make([]slack.Block, 0)

	if imageURL != "" {
		imgTextBlock := slack.NewTextBlockObject("plain_text", "chart", false, false)
		imgBlock := slack.NewImageBlock(imageURL, "chart", "", imgTextBlock)
		blocks = append(blocks, imgBlock)
	}

	if alertText != "" {
		opts = append(opts, slack.MsgOptionAsUser(true), slack.MsgOptionText(alertText, false))
		mainTextBlock := slack.NewTextBlockObject("mrkdwn", alertText, false, false)
		mainSectionBlock := slack.NewSectionBlock(mainTextBlock, nil, nil)
		blocks = append(blocks, mainSectionBlock)
	}

	opts = append(opts, slack.MsgOptionBlocks(blocks...))

	return opts
}
