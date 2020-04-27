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
	ServerName   string `json:"serverName" yaml:"serverName"`
	ServerPort   string `json:"serverPort" yaml:"serverPort"`
	AuthUsername string `json:"authUsername" yaml:"authUsername"`
	AuthPassword string `json:"authPassword" yaml:"authPassword"`
	AuthIdentity string `json:"authIdentity" yaml:"authIdentity"`
	AuthSecret   string `json:"authSecret" yaml:"authSecret"`
	RequireTLS   bool   `json:"requireTLS" yaml:"requireTLS"`
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
		return fmt.Errorf("serverName must be not empty")
	}

	if strings.TrimSpace(cfg.ServerPort) == "" {
		return fmt.Errorf("serverPort must be not empty")
	}

	return nil
}
