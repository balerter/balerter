package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"time"
)

type APIAlertInfo struct {
	Name       string      `json:"name"`
	Level      alert.Level `json:"level"`
	LastChange time.Time   `json:"last_change"`
	Start      time.Time   `json:"start"`
	Count      int         `json:"count"`
}

func (m *Manager) GetAlerts() []*APIAlertInfo {
	m.alertsMx.RLock()
	defer m.alertsMx.RUnlock()

	result := make([]*APIAlertInfo, 0)

	for name, a := range m.alerts {
		alertInfo := &APIAlertInfo{
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
