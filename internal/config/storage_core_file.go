package config

import (
	"fmt"
	"strings"
	"time"
)

const (
	defaultStorageCoreFileTimeout = time.Second
)

type StorageCoreFile struct {
	Name    string        `json:"name" yaml:"name"`
	Path    string        `json:"path" yaml:"path"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func (cfg StorageCoreFile) SetDefaults() {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultStorageCoreFileTimeout
	}
}

func (cfg StorageCoreFile) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.Path) == "" {
		return fmt.Errorf("path must be not empty")
	}

	return nil
}
