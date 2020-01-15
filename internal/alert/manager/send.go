package manager

import (
	"go.uber.org/zap"
)

func (m *Manager) Send(message string, channels []string) error {

	for _, channelName := range channels {
		ch, ok := m.channels[channelName]
		if !ok {
			m.logger.Warn("channel not found", zap.String("channel", channelName))
			continue
		}

		if err := ch.Send("", message); err != nil {
			m.logger.Error("error send message", zap.String("channel", channelName), zap.Error(err))
		}
	}

	return nil
}
