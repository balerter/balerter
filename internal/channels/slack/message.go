package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

func createSlackMessageOptions(alertText, imageURL string, fields map[string]string, level string) []slack.MsgOption {
	opts := make([]slack.MsgOption, 0)
	blocks := make([]slack.Block, 0)

	if imageURL != "" {
		imgTextBlock := slack.NewTextBlockObject("plain_text", "chart", false, false)
		imgBlock := slack.NewImageBlock(imageURL, "chart", "", imgTextBlock)
		blocks = append(blocks, imgBlock)
	}

	if alertText != "" {
		opts = append(opts, slack.MsgOptionAsUser(true))
		mainTextBlock := slack.NewTextBlockObject("mrkdwn", alertText, false, false)
		mainSectionBlock := slack.NewSectionBlock(mainTextBlock, nil, nil)
		blocks = append(blocks, mainSectionBlock)
	}

	fieldsBlocks := make([][]*slack.TextBlockObject, 0)

	ff := make([]*slack.TextBlockObject, 0)

	// Slack supports up to 10 fields
	var idx int
	for key, value := range fields {
		ff = append(ff, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%s = %s", key, value), false, false))

		if (idx+1)%10 == 0 {
			fieldsBlocks = append(fieldsBlocks, ff)
			ff = make([]*slack.TextBlockObject, 0)
		}
		idx++
	}

	if len(fields) > 0 && len(fields)%10 != 0 {
		fieldsBlocks = append(fieldsBlocks, ff)
	}

	for _, fieldBlock := range fieldsBlocks {
		blocks = append(blocks, slack.NewSectionBlock(nil, fieldBlock, nil))
	}

	attachment := slack.Attachment{
		Color:    getColorByLevel(level),
		Fallback: alertText,
		Blocks:   slack.Blocks{BlockSet: blocks},
	}

	opts = append(opts, slack.MsgOptionAttachments(attachment))

	return opts
}

func getColorByLevel(l string) string {
	switch l {
	case "success":
		return "#00aa00"
	case "warning":
		return "#ffcc00"
	case "error":
		return "#ff0000"
	}

	return "#cccccc"
}
