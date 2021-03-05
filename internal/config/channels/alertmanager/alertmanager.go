package alertmanager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"strings"
)

// Alertmanager channel config
type Alertmanager struct {
	// Name of the channel
	Name string `json:"name" yaml:"name"`
	// Settings contains webhook settings
	Settings webhook.Settings `json:"settings" yaml:"settings"`
}

// Validate config
func (cfg Alertmanager) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if err := cfg.Settings.Validate(); err != nil {
		return err
	}

	return nil
}
