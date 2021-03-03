package sqlite

import (
	"fmt"
	"strings"
	"time"
)

type Tables struct {
	Alerts string `json:"alerts" yaml:"alerts" hcl:"alerts"`
	KV     string `json:"kv" yaml:"kv" hcl:"kv"`
}

type Sqlite struct {
	Name    string        `json:"name" yaml:"name" hcl:"name,label"`
	Path    string        `json:"path" yaml:"path" hcl:"path"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
	Tables  Tables        `json:"tables" yaml:"tables" hcl:"tables,block"`
}

func (cfg Sqlite) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Path) == "" {
		return fmt.Errorf("path must be not empty")
	}
	if strings.TrimSpace(cfg.Tables.Alerts) == "" {
		return fmt.Errorf("table Alerts must be not empty")
	}
	if strings.TrimSpace(cfg.Tables.KV) == "" {
		return fmt.Errorf("table KV must be not empty")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
