package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/corestorage/provider/memory"
	"github.com/balerter/balerter/internal/corestorage/provider/sql"
	"go.uber.org/zap"
	"time"
)

type Manager struct {
	storages map[string]coreStorage.CoreStorage
	logger   *zap.Logger
}

func New(cfg *core.Core, logger *zap.Logger) (*Manager, error) {
	m := &Manager{
		storages: map[string]coreStorage.CoreStorage{},
		logger:   logger,
	}

	m.storages["memory"] = memory.New()

	if cfg == nil {
		return m, nil
	}

	for _, c := range cfg.Sqlite {
		s, err := sql.New("sqlite."+c.Name, "sqlite3", c.Path, c.Tables.Alerts, c.Tables.KV, time.Millisecond*time.Duration(c.Timeout), logger)
		if err != nil {
			return nil, fmt.Errorf("error create file storage, %w", err)
		}

		m.storages[s.Name()] = s
	}

	for _, c := range cfg.Postgres {
		connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database,
			c.SSLMode,
			c.SSLCertPath,
		)

		s, err := sql.New("postgres."+c.Name, "postgres", connectionString, c.Tables.Alerts, c.Tables.KV, time.Millisecond*time.Duration(c.Timeout), logger)
		if err != nil {
			return nil, fmt.Errorf("error create postgres storage, %w", err)
		}

		m.storages[s.Name()] = s
	}

	return m, nil
}

func (m *Manager) Get(name string) (coreStorage.CoreStorage, error) {
	if name == "" {
		name = "memory"
	}

	s, ok := m.storages[name]
	if !ok {
		return nil, fmt.Errorf("storage not found")
	}

	return s, nil
}

func (m *Manager) Stop() {
	var err error
	for name, s := range m.storages {
		m.logger.Debug("stop core storage", zap.String("name", name))
		err = s.Stop()
		if err != nil {
			m.logger.Error("error stop core storage", zap.String("name", name), zap.Error(err))
		}
	}
}
