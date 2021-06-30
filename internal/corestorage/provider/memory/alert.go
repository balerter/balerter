package memory

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/metrics"
	"time"
)

func (m *storageAlert) Get(name string) (*alert.Alert, error) {
	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	a, ok := m.alerts[name]
	if !ok {
		return nil, fmt.Errorf("alert not found")
	}

	return a, nil
}

func (m *storageAlert) Index(l []alert.Level) (alert.Alerts, error) {
	var result alert.Alerts
	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	for _, a := range m.alerts {
		for _, l := range l {
			if a.Level == l {
				result = append(result, a)
				break
			}
		}
	}

	return result, nil
}

func (m *storageAlert) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	metrics.SetAlertLevel(name, level)

	m.mxAlerts.Lock()
	defer m.mxAlerts.Unlock()

	a, ok := m.alerts[name]
	if !ok {
		a = alert.New(name)
		a.Level = level
		m.alerts[name] = a
		return a, level != alert.LevelSuccess, nil
	}

	if a.Level == level {
		a.Count++
		return a, false, nil
	}

	a.Level = level
	a.LastChange = time.Now()

	return a, true, nil
}
