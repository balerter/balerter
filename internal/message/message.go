package message

// Message represents a Message struct
type Message struct {
	Level     string            `json:"level"`
	AlertName string            `json:"alert_name"`
	Group     string            `json:"group,omitempty"`
	Text      string            `json:"text"`
	Image     string            `json:"image,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// New returns new Message instance
func New(level, alertName, text string, image string, fields map[string]string, group string) *Message {
	m := &Message{
		Level:     level,
		AlertName: alertName,
		Group:     group,
		Text:      text,
		Image:     image,
		Fields:    fields,
	}

	return m
}
