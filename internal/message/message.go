package message

// Message represents a Message struct
type Message struct {
	Level     string `json:"level"`
	AlertName string `json:"alert_name"`
	Text      string `json:"text"`
	Image     string `json:"image,omitempty"`
}

// New returns new Message instance
func New(level, alertName, text string, image string) *Message {
	m := &Message{
		Level:     level,
		AlertName: alertName,
		Text:      text,
		Image:     image,
	}

	return m
}
