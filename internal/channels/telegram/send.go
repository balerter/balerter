package telegram

import (
	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
)

// Send the message to the channel
func (tg *Telegram) Send(mes *message.Message) error {
	tg.logger.Debug("tg send message")

	if mes.Image != "" {
		tgMessage := api.NewPhotoMessage(tg.chatID, mes.Image, "")
		err := tg.api.SendPhotoMessage(tgMessage)
		if err != nil {
			tg.logger.Error("error send photo", zap.Error(err))
		}
	}

	tgMessage := api.NewTextMessage(tg.chatID, mes.Text)
	return tg.api.SendTextMessage(tgMessage)
}
