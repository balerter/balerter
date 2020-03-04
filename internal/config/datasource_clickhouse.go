package config

import (
	"fmt"
	"strings"
)

type DataSourceClickhouse struct {
	Name        string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Database    string `json:"database" yaml:"database"`
	SSLCertPath string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
}

func (cfg DataSourceClickhouse) SetDefaults() {
	cfg.Host = "127.0.0.1"
	cfg.Port = 9000
	cfg.Username = "default"
	cfg.Database = "default"
}

func (cfg DataSourceClickhouse) Validate() error {
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
