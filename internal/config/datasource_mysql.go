package config

import (
	"fmt"
	"strings"
	"time"
)

type DataSourceMysql struct {
	Name    string        `json:"name" yaml:"name"`
	DSN     string        `json:"dsn" yaml:"dsn"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func (cfg *DataSourceMysql) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.DSN) == "" {
		return fmt.Errorf("DSN must be not empty")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
