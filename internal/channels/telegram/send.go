package telegram

import (
	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/message"
)

func (tg *Telegram) Send(mes *message.Message) error {
	tg.logger.Debug("tg send message")

	if mes.Image != "" {
		tgMessage := api.NewPhotoMessage(tg.chatID, mes.Image, "")
		tg.api.SendPhotoMessage(tgMessage)
	}

	tgMessage := api.NewTextMessage(tg.chatID, mes.Text)
	tg.api.SendTextMessage(tgMessage)

	return nil
}
