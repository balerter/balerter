package mysql

import (
	"fmt"
	"strings"
	"time"
)

type Mysql struct {
	Name    string        `json:"name" yaml:"name" hcl:"name,label"`
	DSN     string        `json:"dsn" yaml:"dsn" hcl:"dsn"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

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
