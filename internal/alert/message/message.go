package message

type Level int

const (
	LevelError Level = iota
	LevelWarn
	LevelSuccess
	LevelInfo
)

type Message struct {
	AlertName string
	Text      string
	Fields    []string
}

func New(alertName, text string, fields ...string) *Message {
	m := &Message{
		AlertName: alertName,
		Text:      text,
		Fields:    fields,
	}

	return m
}
