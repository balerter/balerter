package message

type Message struct {
	Level     string   `json:"level"`
	AlertName string   `json:"alert_name"`
	Text      string   `json:"text"`
	Fields    []string `json:"fields,omitempty"`
	Image     string   `json:"image,omitempty"`
}

func New(level, alertName, text string, fields []string, image string) *Message {
	m := &Message{
		Level:     level,
		AlertName: alertName,
		Text:      text,
		Fields:    fields,
		Image:     image,
	}

	return m
}
