package datasources

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/config/datasources/loki"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/balerter/balerter/internal/util"
)

type DataSources struct {
	Clickhouse []*clickhouse.Clickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []*prometheus.Prometheus `json:"prometheus" yaml:"prometheus"`
	Postgres   []*postgres.Postgres     `json:"postgres" yaml:"postgres"`
	MySQL      []*mysql.Mysql           `json:"mysql" yaml:"mysql"`
	Loki       []*loki.Loki             `json:"loki" yaml:"loki"`
}

func (cfg *DataSources) Validate() error {
	var names []string

	for _, c := range cfg.Clickhouse {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'clickhouse': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Prometheus {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'prometheus': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Postgres {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'postgres': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.MySQL {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'mysql': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Loki {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'loki': %s", name)
	}

	return nil
}
