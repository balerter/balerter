package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/datasource/provider/clickhouse"
	"github.com/balerter/balerter/internal/datasource/provider/loki"
	"github.com/balerter/balerter/internal/datasource/provider/mysql"
	"github.com/balerter/balerter/internal/datasource/provider/postgres"
	"github.com/balerter/balerter/internal/datasource/provider/prometheus"
	moduleMock "github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	"go.uber.org/zap"
)

type Manager struct {
	logger  *zap.Logger
	modules map[string]modules.ModuleTest
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:  logger,
		modules: make(map[string]modules.ModuleTest),
	}

	return m
}

func (m *Manager) Clean() {
	for _, m := range m.modules {
		m.Clean()
	}
}

func (m *Manager) Result() ([]modules.TestResult, error) {
	var result []modules.TestResult
	for _, m := range m.modules {
		results, err := m.Result()
		if err != nil {
			return nil, err
		}
		for _, r := range results {
			r.ModuleName = "datasource." + r.ModuleName
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *Manager) Init(cfg config.DataSources) error {
	for _, clickhouseCfg := range cfg.Clickhouse {
		_ = clickhouseCfg
		mod := moduleMock.New(clickhouse.ModuleName(clickhouseCfg.Name), clickhouse.Methods(), m.logger)
		m.modules[mod.Name()] = mod
	}

	for _, prometheusCfg := range cfg.Prometheus {
		mod := moduleMock.New(prometheus.ModuleName(prometheusCfg.Name), prometheus.Methods(), m.logger)
		m.modules[mod.Name()] = mod
	}

	for _, postgresCfg := range cfg.Postgres {
		mod := moduleMock.New(postgres.ModuleName(postgresCfg.Name), postgres.Methods(), m.logger)
		m.modules[mod.Name()] = mod
	}

	for _, c := range cfg.MySQL {
		mod := moduleMock.New(mysql.ModuleName(c.Name), mysql.Methods(), m.logger)
		m.modules[mod.Name()] = mod
	}

	for _, c := range cfg.Loki {
		mod := moduleMock.New(loki.ModuleName(c.Name), loki.Methods(), m.logger)
		m.modules[mod.Name()] = mod
	}

	return nil
}

func (m *Manager) Get() []modules.ModuleTest {
	mm := make([]modules.ModuleTest, 0)

	for _, mod := range m.modules {
		mm = append(mm, mod)
	}

	return mm
}
