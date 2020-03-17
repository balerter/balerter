package config

import (
	"fmt"
	"strings"
)

type DataSourceMysql struct {
	Name string `json:"name" yaml:"name"`
	DSN  string `json:"dsn" yaml:"dsn"`
}

func (cfg DataSourceMysql) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.DSN) == "" {
		return fmt.Errorf("DSN must be not empty")
	}

	return nil
}
