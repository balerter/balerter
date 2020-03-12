package config

import (
	"fmt"
	"strings"
)

type ChannelNotify struct {
	Name  string             `json:"name" yaml:"name"`
	Icons ChannelNotifyIcons `json:"icons" yaml:"icons"`
}

type ChannelNotifyIcons struct {
	Success string `json:"success" yaml:"success"`
	Error   string `json:"error" yaml:"error"`
	Warning string `json:"warning" yaml:"warning"`
}

func (cfg ChannelNotify) SetDefaults() {
}

func (cfg ChannelNotify) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
