package config

import (
	"fmt"
	"strings"
)

// ChannelEmail configures notifications via email.
type ChannelEmail struct {
	From         string `json:"from,omitempty" yaml:"from,omitempty"`
	To           string `json:"to,omitempty" yaml:"to,omitempty"`
	ServerName   string `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	AuthUsername string `json:"auth_username,omitempty" yaml:"auth_username,omitempty"`
	AuthPassword string `json:"auth_password,omitempty" yaml:"auth_password,omitempty"`
}

// Validate checks the email configuration.
func (cfg ChannelEmail) Validate() error {
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
