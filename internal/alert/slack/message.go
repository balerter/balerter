package slack

type slackMessage struct {
	Blocks []block `json:"blocks"`
}

type block struct {
	Type    string                 `json:"type"`
	Section *blockSection          `json:"text,omitempty"`
	Context *[]blockContextElement `json:"elements,omitempty"`
}

type blockContextElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type blockSection struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func createSlackMessage(alertName, alertText string) slackMessage {
	s := slackMessage{Blocks: []block{
		{
			Type: "section",
			Section: &blockSection{
				Type: "mrkdwn",
				Text: alertText,
			},
		},
	}}

	if alertName != "" {
		s.Blocks = append(s.Blocks, block{
			Type: "context",
			Context: &[]blockContextElement{
				{
					Type: "mrkdwn",
					Text: "Alert name: " + alertName,
				},
			},
		})
	}

	return s
}
