package telegram

import (
	"fmt"
	"strconv"

	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/message"

	"go.uber.org/zap"
)

// Send the message to the channel
func (tg *Telegram) Send(mes *message.Message) error {
	tg.logger.Debug("tg send message")

	if len(mes.Fields) > 0 {
		mes.Text += addFields(mes.Fields)
	}

	if mes.Image != "" {
		tgMessage := api.NewPhotoMessage(tg.chatID, mes.Image, mes.Text)
		err := tg.api.SendPhotoMessage(tgMessage)
		if err != nil {
			tg.logger.Error("error send photo", zap.Error(err))
		}
		return nil
	}

	tgMessage := api.NewTextMessage(tg.chatID, mes.Text)
	return tg.api.SendTextMessage(tgMessage)
}

func addFields(fields map[string]string) string {
	m := strconv.Itoa(maxKeyLen(fields))

	s := "\n\n```\n"
	for k, v := range fields {
		s += fmt.Sprintf("%-"+m+"s = %s\n", k, v)
	}
	s += "\n```"
	return s
}

func maxKeyLen(fields map[string]string) int {
	max := 0
	for k := range fields {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}
