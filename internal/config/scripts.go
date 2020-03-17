package config

import (
	"fmt"
	"time"
)

type Scripts struct {
	UpdateInterval time.Duration  `json:"updateInterval" yaml:"updateInterval"`
	Sources        ScriptsSources `json:"sources" yaml:"sources"`
}

func (cfg Scripts) Validate() error {
	if cfg.UpdateInterval < 0 {
		return fmt.Errorf("updateInterval must be not less than 0")
	}
	return cfg.Sources.Validate()
}
