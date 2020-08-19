package api

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (api *API) SendPhotoMessage(mes *PhotoMessage) error {
	file, err := mes.getPhotoFile()
	if err != nil {
		return fmt.Errorf("error get photo file, %w", err)
	}

	msg := tgbotapi.NewPhotoUpload(mes.ChatID, file)
	api.sendMessage(msg)
	return nil
}

func (api *API) SendTextMessage(mes *TextMessage) error {
	msg := tgbotapi.NewMessage(mes.ChatID, mes.Text)
	api.sendMessage(msg)
	return nil
}

func (api *API) sendMessage(message tgbotapi.Chattable) {
	api.api.Send(message)
}
