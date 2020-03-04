package config

import (
	"fmt"
	"strings"
)

type DataSourcePostgres struct {
	Name        string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Database    string `json:"database" yaml:"database"`
	SSLMode     string `json:"ssl_mode" yaml:"ssl_mode"`
	SSLCertPath string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
}

func (cfg DataSourcePostgres) SetDefaults() {
	cfg.Host = "127.0.0.1"
	cfg.Port = 5432
	cfg.Username = "postgres"
	cfg.Database = "postgres"
}

func (cfg DataSourcePostgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if cfg.Host == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}

	return nil
}
