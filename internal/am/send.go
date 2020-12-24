package manager

import (
	"errors"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/message"
)

var (
	ErrEmptyChannels = errors.New("empty channels")
)

// Send a message
func (m *Manager) Send(level, alertName, text string, options *alert.Options, errs chan<- error) error {
	chs := make(map[string]alertChannel)

	if len(options.Channels) > 0 {
		for _, channelName := range options.Channels {
			// TODO: race on m.channels
			ch, ok := m.channels[channelName]
			if !ok {
				if errs != nil {
					errs <- fmt.Errorf("channel '%s' not found", channelName)
				}
				continue
			}
			chs[channelName] = ch
		}
	} else {
		chs = m.channels
	}

	if len(chs) == 0 {
		return ErrEmptyChannels
	}

	for name, module := range chs {
		if err := module.Send(message.New(level, alertName, text, options.Fields, options.Image)); err != nil {
			if errs != nil {
				errs <- fmt.Errorf("error send message to channel %s, %v", name, err)
			}
		}
	}

	return nil
}
