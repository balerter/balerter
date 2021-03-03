package scripts

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/scripts/sources"
	"time"
)

type Scripts struct {
	UpdateInterval time.Duration   `json:"updateInterval" yaml:"updateInterval"`
	Sources        sources.Sources `json:"sources" yaml:"sources"`
}

func (cfg Scripts) Validate() error {
	if cfg.UpdateInterval < 0 {
		return fmt.Errorf("updateInterval must be not less than 0")
	}
	if err := cfg.Sources.Validate(); err != nil {
		return err
	}
	return nil
}
