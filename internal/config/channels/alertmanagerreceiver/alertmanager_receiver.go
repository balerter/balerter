package alertmanagerreceiver

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"strings"
)

// AlertmanagerReceiver channel config
type AlertmanagerReceiver struct {
	// Name of the channel
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Settings contains webhook settings
	Settings webhook.Settings `json:"settings" yaml:"settings,block"`
}

// Validate config
func (cfg AlertmanagerReceiver) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if err := cfg.Settings.Validate(); err != nil {
		return err
	}

	return nil
}
