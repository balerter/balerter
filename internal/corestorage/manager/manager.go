package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/corestorage/provider/memory"
	"github.com/balerter/balerter/internal/corestorage/provider/sql"
	"go.uber.org/zap"
)

type Manager struct {
	storages map[string]coreStorage.CoreStorage
}

func New(cfg core.Core, logger *zap.Logger) (*Manager, error) {
	m := &Manager{
		storages: map[string]coreStorage.CoreStorage{},
	}

	m.storages["memory"] = memory.New()

	for _, c := range cfg.File {
		s, err := sql.New(c.Name, "sqlite3", c.Path, c.Tables.Alerts, c.Tables.KV, c.Timeout, logger)
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

		s, err := sql.New(c.Name, "postgres", connectionString, c.Tables.Alerts, c.Tables.KV, c.Timeout, logger)
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
