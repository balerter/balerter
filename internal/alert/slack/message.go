package slack

type slackMessage struct {
	Text   string  `json:"text"`
	Blocks []block `json:"blocks"`
}

type block struct {
	Type    string       `json:"type"`
	Section *blockText   `json:"text,omitempty"`
	Context *[]blockText `json:"elements,omitempty"`
}

type blockText struct {
	Type  string `json:"type"`
	Emoji bool   `json:"emoji,omitempty"`
	Text  string `json:"text"`
}

func createSlackMessage(alertName, alertText string) slackMessage {
	s := slackMessage{
		Text: alertText,
		Blocks: []block{
			{
				Type: "section",
				Section: &blockText{
					Type: "mrkdwn",
					Text: alertText,
				},
			},
		},
	}

	if alertName != "" {
		s.Blocks = append(s.Blocks, block{
			Type: "context",
			Context: &[]blockText{
				{
					Type: "mrkdwn",
					Text: "Alert name: " + alertName,
				},
			},
		})
	}

	return s
}
