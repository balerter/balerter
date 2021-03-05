package clickhouse

import (
	"fmt"
	"strings"
	"time"
)

// Clickhouse datasource config
type Clickhouse struct {
	// Name of the datasource
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Host connection value
	Host string `json:"host" yaml:"host" hcl:"host"`
	// Port connection value
	Port int `json:"port" yaml:"port" hcl:"port"`
	// Username connection value
	Username string `json:"username" yaml:"username" hcl:"username"`
	// Password connection value
	Password string `json:"password" yaml:"password" hcl:"password,optional"`
	// Database connection value
	Database string `json:"database" yaml:"database" hcl:"database"`
	// SSLCertPath connection value
	SSLCertPath string `json:"sslCertPath" yaml:"sslCertPath" hcl:"sslCertPath,optional"`
	// Timeout connection value
	Timeout time.Duration `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

// Validate config
func (cfg Clickhouse) Validate() error {
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
