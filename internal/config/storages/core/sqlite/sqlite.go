package sqlite

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"strings"
)

// Sqlite core storage config
type Sqlite struct {
	// Name of the core storage
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Path to database file
	Path string `json:"path" yaml:"path" hcl:"path"`
	// Timeout value
	Timeout int `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`

	// TableAlerts is config for Alerts table
	TableAlerts tables.TableAlerts `json:"tableAlerts" yaml:"tableAlerts" hcl:"tableAlerts,block"`
	// TableKV is config for KV table
	TableKV tables.TableKV `json:"tableKV" yaml:"tableKV" hcl:"tableKV,block"`
}

// Validate config
func (cfg Sqlite) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Path) == "" {
		return fmt.Errorf("path must be not empty")
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
