package manager

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
)

// Send a message
func (m *ChannelsManager) Send(a *alert.Alert, text string, options *alert.Options) {
	if options.Quiet {
		return
	}

	chs := make(map[string]alertChannel)

	if len(options.Channels) > 0 {
		for _, channelName := range options.Channels {
			ch, ok := m.channels[channelName]
			if !ok {
				m.logger.Warn("channel not found", zap.String("name", channelName))
				continue
			}
			chs[channelName] = ch
		}
	} else {
		for _, ch := range m.channels {
			if !ch.Ignore() {
				chs[ch.Name()] = ch
			}
		}
	}

	if len(chs) == 0 {
		m.logger.Warn("the message was not sent, empty channels")
		return
	}

	for name, module := range chs {
		if err := module.Send(message.New(a.Level.String(), a.Name, text, options.Image, options.Fields)); err != nil {
			m.logger.Error("error send the message to the channel", zap.String("channel name", name), zap.Error(err))
		}
	}
}
