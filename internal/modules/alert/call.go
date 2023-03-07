package alert

import (
	"github.com/balerter/balerter/internal/alert"
)

func (a *Alert) call(name, text string, scriptChannels []string, escalate map[int][]string, alertLevel alert.Level, options *alert.Options) (*alert.Alert, bool, error) {
	if len(options.Channels) == 0 {
		options.Channels = scriptChannels
	}

	updatedAlert, levelWasUpdated, err := a.storage.Update(name, alertLevel)
	if err != nil {
		return nil, false, err
	}

	// For current Error level check if we need to escalate
	if updatedAlert.Level == alert.LevelError {
		for num, channels := range escalate {
			if updatedAlert.Count == num {
				opts := options
				opts.Channels = channels
				a.chManager.Send(updatedAlert, text, opts)
			}
		}
	}

	if levelWasUpdated || (options.Repeat > 0 && updatedAlert.Count%options.Repeat == 0) {
		a.chManager.Send(updatedAlert, text, options)
	}

	return updatedAlert, levelWasUpdated, nil
}
