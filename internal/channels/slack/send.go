package slack

import (
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
)

func (m *Slack) Send(mes *message.Message) error {
	opts := createSlackMessageOptions(mes.Text, mes.Image)

	_channel, _timestamp, _text, err := m.api.SendMessage(m.channel, opts...)

	m.logger.Debug("send slack message",
		zap.String("channel", _channel),
		zap.String("timestamp", _timestamp),
		zap.String("text", _text),
		zap.Error(err),
	)

	return err
}
