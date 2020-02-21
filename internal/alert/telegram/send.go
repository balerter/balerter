package telegram

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (tg *Telegram) Send(level alert.Level, mes *message.Message, chartData *chartModule.Data) error {

	//log.Printf("Authorized on account %s", tg.bot.Self.UserName)

	//if mes.Image != "" {
	//	_, err := tg.bot.SetChatPhoto(tgbotapi.NewSetChatPhotoUpload(tg.chatID, mes.Image))
	//	if err != nil {
	//		tg.logger.Error("error send image to telegram", zap.Error(err))
	//	}
	//}

	msg := tgbotapi.NewMessage(tg.chatID, mes.Text)
	_, err := tg.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
