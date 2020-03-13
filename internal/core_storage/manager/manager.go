package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	coreStorage "github.com/balerter/balerter/internal/core_storage"
	"github.com/balerter/balerter/internal/core_storage/provider/file"
	"github.com/balerter/balerter/internal/core_storage/provider/memory"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	mx       sync.RWMutex
	storages map[string]coreStorage.CoreStorage
}

func New(cfg config.StoragesCore, logger *zap.Logger) (*Manager, error) {
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
	s, ok := m.storages[name]
	if !ok {
		return nil, fmt.Errorf("storage not found")
	}

	return s, nil
}
