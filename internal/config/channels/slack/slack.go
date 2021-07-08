package slack

import (
	"fmt"
	"strings"
)

// Slack channel config
type Slack struct {
	// Name of the channel
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Token is auth token for slack app
	Token string `json:"token" yaml:"token" hcl:"token"`
	// Channel name
	Channel string `json:"channel" yaml:"channel"  hcl:"channel"`
}

// Validate config
func (cfg Slack) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("token must be not empty")
	}
	if strings.TrimSpace(cfg.Channel) == "" {
		return fmt.Errorf("channel must be not empty")
	}

	return nil
}
