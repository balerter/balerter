package config

import "fmt"

type DataSources struct {
	Clickhouse []DataSourceClickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []DataSourcePrometheus `json:"prometheus" yaml:"prometheus"`
	Postgres   []DataSourcePostgres   `json:"postgres" yaml:"postgres"`
	MySQL      []DataSourceMysql      `json:"mysql" yaml:"mysql"`
	Loki       []DataSourceLoki       `json:"loki" yaml:"loki"`
}

func (cfg DataSources) Validate() error {
	var names []string

	for _, c := range cfg.Clickhouse {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'clickhouse': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Prometheus {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'prometheus': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Postgres {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'postgres': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.MySQL {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'mysql': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Loki {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for datasource 'loki': %s", name)
	}

	return nil
}
