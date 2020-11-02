package config

import (
	"fmt"
	"strings"
)

type ChannelAlertmanagerReceiver struct {
	Name     string           `json:"name" yaml:"name"`
	Settings *WebhookSettings `json:"settings" yaml:"settings"`
}

func (cfg *ChannelAlertmanagerReceiver) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return cfg.Settings.Validate()
}
