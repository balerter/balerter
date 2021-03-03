package alertmanager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"strings"
)

type Alertmanager struct {
	Name     string           `json:"name" yaml:"name"`
	Settings webhook.Settings `json:"settings" yaml:"settings"`
}

func (cfg Alertmanager) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if err := cfg.Settings.Validate(); err != nil {
		return err
	}

	return nil
}
