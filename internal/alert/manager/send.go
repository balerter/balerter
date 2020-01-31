package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"go.uber.org/zap"
)

func (m *Manager) Send(level alert.Level, alertName, text string, channels []string, fields []string) {
	for name, module := range m.channels {
		if len(channels) > 0 && !contains(name, channels) {
			m.logger.Debug("skip send message to channel", zap.String("channel name", name))
			continue
		}

		if err := module.Send(level, message.New(alertName, text, fields)); err != nil {
			m.logger.Error("error send message to channel", zap.String("channel name", name), zap.Error(err))
		}
	}
}

func contains(v string, arr []string) bool {
	for _, item := range arr {
		if item == v {
			return true
		}
	}
	return false
}
