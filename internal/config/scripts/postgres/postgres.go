package postgres

import (
	"fmt"
	"strings"
)

type Fields struct {
	Name string `json:"name" yaml:"name" hcl:"name"`
	Body string `json:"body" yaml:"body" hcl:"body"`
}

type Postgres struct {
	Name        string `json:"name" yaml:"name" hcl:"name,label"`
	Host        string `json:"host" yaml:"host" hcl:"host"`
	Port        int    `json:"port" yaml:"port" hcl:"port"`
	Username    string `json:"username" yaml:"username" hcl:"username"`
	Password    string `json:"password" yaml:"password" hcl:"password"`
	Database    string `json:"database" yaml:"database" hcl:"database"`
	SSLMode     string `json:"sslMode" yaml:"sslMode" hcl:"sslMode,optional"`
	SSLCertPath string `json:"sslCertPath" yaml:"sslCertPath" hcl:"sslCertPath,optional"`
	Timeout     int    `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
	Table       string `json:"table" yaml:"table" hcl:"table"`
	Fields      Fields `json:"fields" yaml:"fields" hcl:"fields,block"`
}

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
	if cfg.Table == "" {
		return fmt.Errorf("table must be defined")
	}
	if err := cfg.Fields.Validate(); err != nil {
		return err
	}

	return nil
}

func (f Fields) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("field name must be defined")
	}
	if f.Body == "" {
		return fmt.Errorf("field body must be defined")
	}
	return nil
}
