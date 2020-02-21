package api

type TextMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type PhotoMessage struct {
	ChatID  int64  `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption,omitempty"`
}

func NewPhotoMessage(chatID int64, photo, text string) *PhotoMessage {
	return &PhotoMessage{
		ChatID:  chatID,
		Photo:   photo,
		Caption: text,
	}
}

func NewTextMessage(chatID int64, text string) *TextMessage {
	return &TextMessage{
		ChatID: chatID,
		Text:   text,
	}
}
