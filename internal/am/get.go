package manager

import (
	"errors"
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"strings"
)

var (
	ErrEmptyName = errors.New("empty alert name")
)

func (m *Manager) Get(name string) (*alert2.Alert, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrEmptyName
	}

	a, err := m.storage.Alert().Get(name)
	if err != nil {
		return nil, fmt.Errorf("error get alert, %w", err)
	}

	return a, nil
}
