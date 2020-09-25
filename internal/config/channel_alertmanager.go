package config

import (
	"fmt"
	"strings"
)

type ChannelAlertmanager struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	URL     string `json:"url" yaml:"url"`
}

func (cfg *ChannelAlertmanager) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	ver := strings.TrimSpace(cfg.Version)
	if ver != "" && ver != "v1" && ver != "v2" {
		return fmt.Errorf("version must be empty or v1 or v2")
	}

	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("url must be not empty")
	}

	return nil
}
