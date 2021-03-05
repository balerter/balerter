package notify

import (
	"fmt"
	"strings"
)

// Notify config
type Notify struct {
	// Name of the channel
	Name string `json:"name" yaml:"name"`
	// Icons settings
	Icons ChannelNotifyIcons `json:"icons" yaml:"icons"`
}

// ChannelNotifyIcons is icon settings
type ChannelNotifyIcons struct {
	// Success icon
	Success string `json:"success" yaml:"success"`
	// Error icon
	Error string `json:"error" yaml:"error"`
	// Warning icon
	Warning string `json:"warning" yaml:"warning"`
}

// Validate config
func (cfg Notify) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
