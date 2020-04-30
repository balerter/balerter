package config

import (
	"fmt"
	"strings"
)

// ChannelEmail configures notifications via email.
type ChannelEmail struct {
	Name       string `json:"name" yaml:"name"`
	From       string `json:"from" yaml:"from"`
	To         string `json:"to" yaml:"to"`
	Cc         string `json:"cc" yaml:"cc"`
	Host       string `json:"host" yaml:"host"`
	Port       string `json:"port" yaml:"port"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	WithoutTLS bool   `json:"withoutTLS" yaml:"withoutTLS"`
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

	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be not empty")
	}

	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("port must be not empty")
	}

	return nil
}
