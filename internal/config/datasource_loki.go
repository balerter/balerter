package config

import (
	"fmt"
	"strings"
)

type DataSourceLoki struct {
	Name      string    `json:"name" yaml:"name"`
	URL       string    `json:"url" yaml:"url"`
	BasicAuth BasicAuth `json:"basicAuth" yaml:"basicAuth"`
}

func (cfg DataSourceLoki) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("url must be not empty")
	}

	return nil
}
