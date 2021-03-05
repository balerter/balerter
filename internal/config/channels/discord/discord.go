package discord

import (
	"fmt"
	"strings"
)

// Discord channel config
type Discord struct {
	// Name of the channel
	Name string `json:"name" yaml:"name"`
	// Token is auth token
	Token string `json:"token" yaml:"token"`
	// ChannelID of a discord channel
	ChannelID int64 `json:"channelId" yaml:"channelId"`
}

// Validate config
func (cfg Discord) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("token must be not empty")
	}

	if cfg.ChannelID < 1 {
		return fmt.Errorf("channel id must be not empty")
	}

	return nil
}
