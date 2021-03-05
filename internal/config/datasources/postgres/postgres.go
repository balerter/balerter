package postgres

import (
	"fmt"
	"strings"
)

// Postgres datasource config
type Postgres struct {
	// Name of the datasource
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Host value
	Host string `json:"host" yaml:"host" hcl:"host"`
	// Port value
	Port int `json:"port" yaml:"port" hcl:"port"`
	// Username value
	Username string `json:"username" yaml:"username" hcl:"username"`
	// Password value
	Password string `json:"password" yaml:"password" hcl:"password"`
	// Database value
	Database string `json:"database" yaml:"database" hcl:"database"`
	// SSLMode value
	SSLMode string `json:"sslMode" yaml:"sslMode" hcl:"sslMode,optional"`
	// SSLCertPath value
	SSLCertPath string `json:"sslCertPath" yaml:"sslCertPath" hcl:"sslCertPath,optional"`
	// Timeout value
	Timeout int `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

// Validate config
func (cfg Postgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if cfg.Host == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
