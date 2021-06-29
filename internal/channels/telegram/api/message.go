package api

// TextMessage represents TextMessage from Telegram API
type TextMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// PhotoMessage represents PhotoMessage from Telegram API
type PhotoMessage struct {
	ChatID  int64  `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption,omitempty"`
}

// NewPhotoMessage creates new PhotoMessage
func NewPhotoMessage(chatID int64, photo, text string) *PhotoMessage {
	return &PhotoMessage{
		ChatID:  chatID,
		Photo:   photo,
		Caption: text,
	}
}

// NewTextMessage creates new TextMessage
func NewTextMessage(chatID int64, text string) *TextMessage {
	return &TextMessage{
		ChatID: chatID,
		Text:   text,
	}
}
