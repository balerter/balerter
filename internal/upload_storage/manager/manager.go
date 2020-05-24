package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/upload_storage/provider/s3"
	"go.uber.org/zap"
)

type Provider interface {
}

type Manager struct {
	logger *zap.Logger

	modules map[string]modules.Module
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:  logger,
		modules: make(map[string]modules.Module),
	}

	return m
}

func (m *Manager) Init(cfg config.StoragesUpload) error {
	for idx := range cfg.S3 {
		module, err := s3.New(cfg.S3[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init storage provider s3, %w", err)
		}

		m.modules[module.Name()] = module
	}

	return nil
}

func (m *Manager) Get() []modules.Module {
	mm := make([]modules.Module, 0)

	for _, module := range m.modules {
		mm = append(mm, module)
	}

	return mm
}
