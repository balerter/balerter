package message

type Message struct {
	AlertName string   `json:"alert_name"`
	Text      string   `json:"text"`
	Fields    []string `json:"fields,omitempty"`
	Image     string   `json:"image,omitempty"`
}

func New(alertName, text string, fields []string, image string) *Message {
	m := &Message{
		AlertName: alertName,
		Text:      text,
		Fields:    fields,
		Image:     image,
	}

	return m
}
