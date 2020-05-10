package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/datasource/provider/clickhouse"
	"github.com/balerter/balerter/internal/datasource/provider/loki"
	"github.com/balerter/balerter/internal/datasource/provider/mysql"
	"github.com/balerter/balerter/internal/datasource/provider/postgres"
	"github.com/balerter/balerter/internal/datasource/provider/prometheus"
	"github.com/balerter/balerter/internal/modules"
	"go.uber.org/zap"
)

type Manager struct {
	logger  *zap.Logger
	modules map[string]modules.Module
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:  logger,
		modules: make(map[string]modules.Module),
	}

	return m
}

func (m *Manager) Init(cfg *config.DataSources) error {
	for _, clickhouseCfg := range cfg.Clickhouse {
		module, err := clickhouse.New(&clickhouseCfg, m.logger)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module
	}

	for _, prometheusCfg := range cfg.Prometheus {
		module, err := prometheus.New(prometheusCfg, m.logger)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module
	}

	for _, postgresCfg := range cfg.Postgres {
		module, err := postgres.New(postgresCfg, m.logger)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module
	}

	for _, c := range cfg.MySQL {
		module, err := mysql.New(c, m.logger)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module
	}

	for _, c := range cfg.Loki {
		module, err := loki.New(c, m.logger)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module
	}

	return nil
}

func (m *Manager) Stop() {
	for _, module := range m.modules {
		m.logger.Debug("stop module", zap.String("name", module.Name()))
		err := module.Stop()
		if err != nil {
			m.logger.Error("error stop module", zap.String("name", module.Name()), zap.Error(err))
		}
	}
}

func (m *Manager) Get() []modules.Module {
	mm := make([]modules.Module, 0)

	for _, module := range m.modules {
		mm = append(mm, module)
	}

	return mm
}
