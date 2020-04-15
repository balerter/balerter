package memory

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
)

func (m *storageAlert) GetOrNew(name string) (*alert.Alert, error) {
	m.mxAlerts.RLock()
	a, ok := m.alerts[name]
	m.mxAlerts.RUnlock()
	if ok {
		return a, nil
	}

	m.mxAlerts.Lock()
	defer m.mxAlerts.Unlock()

	a, ok = m.alerts[name]
	if ok {
		return a, nil
	}
	a = alert.AcquireAlert()
	a.SetName(name)
	m.alerts[name] = a
	return a, nil
}

func (m *storageAlert) All() ([]*alert.Alert, error) {
	var result []*alert.Alert

	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	for _, a := range m.alerts {
		result = append(result, a)
	}

	return result, nil
}

func (m *storageAlert) Release(_ *alert.Alert) error {
	return nil
}

func (m *storageAlert) Get(name string) (*alert.Alert, error) {
	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	a, ok := m.alerts[name]
	if !ok {
		return nil, fmt.Errorf("alert not found")
	}

	return a, nil
}
