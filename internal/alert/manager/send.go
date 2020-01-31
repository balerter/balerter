package manager

import (
	"github.com/balerter/balerter/internal/alert/message"
	"go.uber.org/zap"
)

func (m *Manager) Send(level message.Level, alertName, text string, channels map[string]struct{}, fields ...string) {
	for name, module := range m.channels {
		if channels != nil {
			if _, ok := channels[name]; len(channels) > 0 && !ok {
				m.logger.Debug("skip send message to channel", zap.String("channel name", name))
				continue
			}
		}

		if err := module.Send(level, message.New(alertName, text, fields...)); err != nil {
			m.logger.Error("error send message to channel", zap.String("channel name", name), zap.Error(err))
		}
	}
}
