package postgres

import (
	"fmt"
	"strings"
	"time"
)

type Tables struct {
	Alerts string `json:"alerts" yaml:"alerts"`
	KV     string `json:"kv" yaml:"kv"`
}

type Postgres struct {
	Name        string        `json:"name" yaml:"name"`
	Host        string        `json:"host" yaml:"host"`
	Port        int           `json:"port" yaml:"port"`
	Username    string        `json:"username" yaml:"username"`
	Password    string        `json:"password" yaml:"password"`
	Database    string        `json:"database" yaml:"database"`
	SSLMode     string        `json:"sslMode" yaml:"sslMode"`
	SSLCertPath string        `json:"sslCertPath" yaml:"sslCertPath"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`

	Tables Tables `json:"tables" yaml:"tables"`
}

func (cfg Postgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be defined")
	}
	if strings.TrimSpace(cfg.Tables.Alerts) == "" {
		return fmt.Errorf("table Alerts must be not empty")
	}
	if strings.TrimSpace(cfg.Tables.KV) == "" {
		return fmt.Errorf("table KV must be not empty")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
