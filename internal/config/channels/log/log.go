package log

import (
	"fmt"
	"strings"
)

type Log struct {
	Name   string `json:"name" yaml:"name" hcl:"name,label"`
	Ignore bool   `json:"ignore" yaml:"ignore" hcl:"ignore,optional"`
}

// Validate checks the webhook configuration.
func (cfg Log) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
