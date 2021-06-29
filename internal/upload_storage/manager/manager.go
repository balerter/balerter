package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/upload_storage/provider/s3"
	"go.uber.org/zap"
)

// Manager represents the upload storage manager
type Manager struct {
	logger *zap.Logger

	modules map[string]modules.Module
}

// New creates new upload storage manager
func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:  logger,
		modules: make(map[string]modules.Module),
	}

	return m
}

// Init the upload storage manager
func (m *Manager) Init(cfg *upload.Upload) error {
	if cfg == nil {
		return nil
	}
	for idx := range cfg.S3 {
		module, err := s3.New(cfg.S3[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init storage provider s3, %w", err)
		}

		m.modules[module.Name()] = module
	}

	return nil
}

// Get returns the upload storage modules
func (m *Manager) Get() []modules.Module {
	mm := make([]modules.Module, 0)

	for _, module := range m.modules {
		mm = append(mm, module)
	}

	return mm
}
