package config

type DataSources struct {
	Clickhouse []DataSourceClickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []DataSourcePrometheus `json:"prometheus" yaml:"prometheus"`
	Postgres   []DataSourcePostgres   `json:"postgres" yaml:"postgres"`
}

func (cfg DataSources) SetDefaults() {
	for _, c := range cfg.Clickhouse {
		c.SetDefaults()
	}
	for _, c := range cfg.Prometheus {
		c.SetDefaults()
	}
	for _, c := range cfg.Postgres {
		c.SetDefaults()
	}
}

func (cfg DataSources) Validate() error {
	for _, c := range cfg.Clickhouse {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.Prometheus {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.Postgres {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
