package postgres

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"strings"
)

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

	TableAlerts tables.TableAlerts `json:"tableAlerts" yaml:"tableAlerts" hcl:"tableAlerts,block"`
	TableKV     tables.TableKV     `json:"tableKV" yaml:"tableKV" hcl:"tableKV,block"`
}

func (cfg Postgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	if err := cfg.TableAlerts.Validate(); err != nil {
		return err
	}
	if err := cfg.TableKV.Validate(); err != nil {
		return err
	}

	return nil
}
