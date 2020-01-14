package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/modules"
)

type Manager struct {
}

func New() *Manager {
	m := &Manager{}

	return m
}

func (m *Manager) Init(cfg config.DataSources) error {
	return nil
}

func (m *Manager) Get() []modules.Module {
	return nil
}
