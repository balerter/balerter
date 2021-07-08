package notify

import (
	"fmt"
	"strings"
)

// Notify config
type Notify struct {
	// Name of the channel
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Icons settings
	Icons ChannelNotifyIcons `json:"icons" yaml:"icons" hcl:"icons,block"`
}

// ChannelNotifyIcons is icon settings
type ChannelNotifyIcons struct {
	// Success icon
	Success string `json:"success" yaml:"success" hcl:"success,optional"`
	// Error icon
	Error string `json:"error" yaml:"error" hcl:"success,optional"`
	// Warning icon
	Warning string `json:"warning" yaml:"warning" hcl:"success,optional"`
}

// Validate config
func (cfg Notify) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
