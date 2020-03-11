package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"time"
)

type AlertInfo struct {
	Name       string      `json:"name"`
	Level      alert.Level `json:"level"`
	LastChange time.Time   `json:"last_change"`
	Start      time.Time   `json:"start"`
	Count      int         `json:"count"`
}

func (m *Manager) GetAlerts() []*AlertInfo {
	m.alertsMx.RLock()
	defer m.alertsMx.RUnlock()

	result := make([]*AlertInfo, 0)

	for name, a := range m.alerts {
		alertInfo := &AlertInfo{
			Name:       name,
			Level:      a.Level(),
			LastChange: a.GetLastChangeTime(),
			Start:      a.GetStartTime(),
			Count:      a.Count(),
		}

		result = append(result, alertInfo)
	}

	return result
}
