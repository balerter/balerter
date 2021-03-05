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

// DataSources config
type DataSources struct {
	// Clickhouse for clickhouse data sources
	Clickhouse []clickhouse.Clickhouse `json:"clickhouse" yaml:"clickhouse" hcl:"clickhouse,block"`
	// Prometheus for prometheus-like data sources
	Prometheus []prometheus.Prometheus `json:"prometheus" yaml:"prometheus" hcl:"prometheus,block"`
	// Postgres for postgres data sources
	Postgres []postgres.Postgres `json:"postgres" yaml:"postgres" hcl:"postgres,block"`
	// MySQL for mysql data sources
	MySQL []mysql.Mysql `json:"mysql" yaml:"mysql" hcl:"mysql,block"`
	// Loki for Loki data sources
	Loki []loki.Loki `json:"loki" yaml:"loki" hcl:"loki,block"`
}

// Validate config
func (cfg DataSources) Validate() error {
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
