package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/datasource/provider/clickhouse"
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

func (m *Manager) Init(cfg config.DataSources) error {

	for _, clickhouseCfg := range cfg.Clickhouse {
		module, err := clickhouse.New(clickhouseCfg)
		if err != nil {
			return err
		}
		m.modules[module.Name()] = module

	}

	return nil
}

func (m *Manager) Get() []modules.Module {
	return nil
}
