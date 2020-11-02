package config

import (
	"fmt"
	"strings"
)

type ChannelAlertmanager struct {
	Name     string           `json:"name" yaml:"name"`
	Settings *WebhookSettings `json:"settings"`
}

func (cfg *ChannelAlertmanager) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if cfg.Settings == nil {
		return fmt.Errorf("sttings must be defined")
	}

	return cfg.Settings.Validate()
}
