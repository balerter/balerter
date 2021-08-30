package memory

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/metrics"
	"time"
)

func (m *storageAlert) Get(name string) (*alert.Alert, error) {
	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	a, ok := m.alerts[name]
	if !ok {
		return nil, nil
	}

	return a, nil
}

func (m *storageAlert) Index(l []alert.Level) (alert.Alerts, error) {
	var result alert.Alerts
	m.mxAlerts.RLock()
	defer m.mxAlerts.RUnlock()

	for _, a := range m.alerts {
		if len(l) > 0 {
			for _, l := range l {
				if a.Level == l {
					result = append(result, a)
					break
				}
			}
			continue
		}
		result = append(result, a)
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

	a.Count = 1
	a.Level = level
	a.LastChange = time.Now()

	return a, true, nil
}
