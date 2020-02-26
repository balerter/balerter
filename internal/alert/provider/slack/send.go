package slack

import (
	"github.com/balerter/balerter/internal/alert/message"
	"go.uber.org/zap"
)

func (m *Slack) Send(message *message.Message) error {

	opts := createSlackMessageOptions(message.Text, message.Image, message.Fields...)

	_channel, _timestamp, _text, err := m.api.SendMessage(m.channel, opts...)

	m.logger.Debug("send slack message", zap.String("channel", _channel), zap.String("timestamp", _timestamp), zap.String("text", _text), zap.Error(err))

	return err
}
