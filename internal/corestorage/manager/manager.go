package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/corestorage/provider/file"
	"github.com/balerter/balerter/internal/corestorage/provider/memory"
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
		s, err := file.New(c, logger)
		if err != nil {
			return nil, fmt.Errorf("error create file storage, %w", err)
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
