package core

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/postgres"
	"github.com/balerter/balerter/internal/config/storages/core/sqlite"
	"github.com/balerter/balerter/internal/util"
)

// Core storages config
type Core struct {
	// Sqlite configs
	Sqlite []sqlite.Sqlite `json:"sqlite" yaml:"sqlite" hcl:"sqlite,block"`
	// Postgres config
	Postgres []postgres.Postgres `json:"postgres" yaml:"postgres" hcl:"postgres,block"`
}

// Validate config
func (cfg Core) Validate() error {
	var names []string

	for _, c := range cfg.Sqlite {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for core storages 'sqlite': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Postgres {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for core storages 'postgres': %s", name)
	}

	return nil
}
