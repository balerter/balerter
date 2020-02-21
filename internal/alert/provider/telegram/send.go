package telegram

import (
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/alert/provider/telegram/api"
	"go.uber.org/zap"
)

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
