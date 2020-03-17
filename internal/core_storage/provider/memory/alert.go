package memory

import (
	"github.com/balerter/balerter/internal/alert/alert"
)

func (m *Memory) GetOrNew(name string) (*alert.Alert, error) {
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

func (m *Memory) All() ([]*alert.Alert, error) {
	var result []*alert.Alert

	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	for _, a := range m.alerts {
		result = append(result, a)
	}

	return result, nil
}

func (m *Memory) Release(_ *alert.Alert) {}
