package manager

import (
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
)

// Send a message
func (m *Manager) Send(level, alertName, text string, channels, fields []string, image string) {
	chs := make(map[string]alertChannel)

	if len(channels) > 0 {
		for _, channelName := range channels {
			// TODO: race on m.channels
			ch, ok := m.channels[channelName]
			if !ok {
				m.logger.Error("channel not found", zap.String("channel name", channelName))
				continue
			}
			chs[channelName] = ch
		}
	} else {
		chs = m.channels
	}

	if len(chs) == 0 {
		m.logger.Error("empty channels")
		return
	}

	for name, module := range chs {
		if err := module.Send(message.New(level, alertName, text, fields, image)); err != nil {
			m.logger.Error("error send message to channel", zap.String("channel name", name), zap.Error(err))
		}
	}
}
