package mysql

import (
	"fmt"
	"strings"
)

// Mysql datasource config
type Mysql struct {
	// Name of the datasource
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// DSN connection data
	DSN string `json:"dsn" yaml:"dsn" hcl:"dsn"`
	// Timeout value
	Timeout int `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

// Validate config
func (cfg Mysql) Validate() error {
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
