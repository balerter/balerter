package manager

import (
	"github.com/balerter/balerter/internal/config/scripts"
	fileProvider "github.com/balerter/balerter/internal/script/provider/file"
	folderProvider "github.com/balerter/balerter/internal/script/provider/folder"
	"github.com/balerter/balerter/internal/script/provider/postgres"
	"github.com/balerter/balerter/internal/script/script"
)

// Provider is an interface for script provider
type Provider interface {
	Get() ([]*script.Script, error)
}

// Manager represents the script manager
type Manager struct {
	providers map[string]Provider
}

// New creates new script manager
func New() *Manager {
	m := &Manager{
		providers: make(map[string]Provider),
	}

	return m
}

// Init the script manager
func (m *Manager) Init(cfg *scripts.Scripts) error {
	if cfg == nil {
		return nil
	}
	for _, folderConfig := range cfg.Folder {
		m.providers[folderConfig.Name] = folderProvider.New(folderConfig)
	}

	for _, c := range cfg.File {
		m.providers[c.Name] = fileProvider.New(c)
	}

	for _, c := range cfg.Postgres {
		p, err := postgres.New(c)
		if err != nil {
			return err
		}
		m.providers[c.Name] = p
	}

	return nil
}

// Get returns scripts obtained from providers
func (m *Manager) Get() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	for _, p := range m.providers {
		s, err := p.Get()
		if err != nil {
			return nil, err
		}

		ss = append(ss, s...)
	}

	ss = removeTests(ss)

	return ss, nil
}

func removeTests(ss []*script.Script) []*script.Script {
	var i int
	for idx, s := range ss {
		if s.IsTest {
			continue
		}
		ss[i] = ss[idx]
		i++
	}
	ss = ss[:i]

	return ss
}

// GetWithTests returns scripts with tests
func (m *Manager) GetWithTests() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	for _, p := range m.providers {
		s, err := p.Get()
		if err != nil {
			return nil, err
		}

		ss = append(ss, s...)
	}

	return ss, nil
}
