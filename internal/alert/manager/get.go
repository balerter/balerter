package manager

import (
	"errors"
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"strings"
)

var (
	ErrEmptyName = errors.New("empty alert name")
)

func (m *Manager) Get(name string) (*alert.Alert, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrEmptyName
	}

	a, err := m.engine.Alert().Get(name)
	if err != nil {
		return nil, fmt.Errorf("error get alert, %w", err)
	}

	return a, nil
}
