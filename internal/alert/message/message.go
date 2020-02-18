package message

type Message struct {
	AlertName string
	Text      string
	Fields    []string
	Image     string
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
