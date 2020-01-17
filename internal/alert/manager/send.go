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

func (m *Manager) sendSuccess(alertName, message string) {
	for name, module := range m.channels {
		if err := module.SendSuccess(alertName, message); err != nil {
			m.logger.Error("error send message to channel", zap.String("name", name), zap.Error(err))
		}
	}
}

func (m *Manager) sendError(alertName, message string) {
	for name, module := range m.channels {
		if err := module.SendError(alertName, message); err != nil {
			m.logger.Error("error send message to channel", zap.String("name", name), zap.Error(err))
		}
	}
}
