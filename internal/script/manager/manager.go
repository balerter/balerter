package manager

import (
	"github.com/balerter/balerter/internal/config/scripts/sources"
	fileProvider "github.com/balerter/balerter/internal/script/provider/file"
	folderProvider "github.com/balerter/balerter/internal/script/provider/folder"
	"github.com/balerter/balerter/internal/script/script"
)

type Provider interface {
	Get() ([]*script.Script, error)
}

type Manager struct {
	providers map[string]Provider
}

func New() *Manager {
	m := &Manager{
		providers: make(map[string]Provider),
	}

	return m
}

func (m *Manager) Init(cfg sources.Sources) error {
	for _, folderConfig := range cfg.Folder {
		m.providers[folderConfig.Name] = folderProvider.New(folderConfig)
	}

	for _, c := range cfg.File {
		m.providers[c.Name] = fileProvider.New(c)
	}

	return nil
}

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
