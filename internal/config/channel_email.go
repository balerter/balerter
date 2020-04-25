package config

import (
	"fmt"
	"strings"
)

// ChannelEmail configures notifications via email.
type ChannelEmail struct {
	Name         string `json:"name" yaml:"name"`
	From         string `json:"from" yaml:"from"`
	To           string `json:"to" yaml:"to"`
	ServerName   string `json:"server_name" yaml:"server_name"`
	AuthUsername string `json:"auth_username" yaml:"auth_username"`
	AuthPassword string `json:"auth_password" yaml:"auth_password"`
	AuthIdentity string `json:"auth_identity" yaml:"auth_identity"`
	AuthSecret   string `json:"auth_secret" yaml:"auth_secret"`
}

// Validate checks the email configuration.
func (cfg ChannelEmail) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.From) == "" {
		return fmt.Errorf("from must be not empty")
	}

	if strings.TrimSpace(cfg.To) == "" {
		return fmt.Errorf("to must be not empty")
	}

	if strings.TrimSpace(cfg.ServerName) == "" {
		return fmt.Errorf("server_name must be not empty")
	}

	return nil
}
