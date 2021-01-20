package file

import (
	"fmt"
	"strings"
	"time"
)

const (
	defaultTablesAlerts = "alerts"
	defaultTablesKV     = "kv"
)

type Tables struct {
	Alerts string `json:"alerts" yaml:"alerts"`
	KV     string `json:"kv" yaml:"kv"`
}

type File struct {
	Name    string        `json:"name" yaml:"name"`
	Path    string        `json:"path" yaml:"path"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Tables  Tables        `json:"tables" yaml:"table"`
}

func (cfg *File) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Path) == "" {
		return fmt.Errorf("path must be not empty")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}

// TODO: call defaults, or use `aconfig` for set defaults
func (cfg *File) Defaults() {
	if cfg.Tables.Alerts == "" {
		cfg.Tables.Alerts = defaultTablesAlerts
	}
	if cfg.Tables.KV == "" {
		cfg.Tables.KV = defaultTablesKV
	}
}
