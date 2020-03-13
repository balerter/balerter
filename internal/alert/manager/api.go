package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"time"
)

type AlertInfo struct {
	Name       string      `json:"name"`
	Level      alert.Level `json:"level"`
	LastChange time.Time   `json:"last_change"`
	Count      int         `json:"count"`
}

func (m *Manager) GetAlerts() []*AlertInfo {

	result := make([]*AlertInfo, 0)

	for _, a := range m.engine.All() {
		alertInfo := &AlertInfo{
			Name:       a.Name(),
			Level:      a.Level(),
			LastChange: a.GetLastChangeTime(),
			Count:      a.Count(),
		}

		result = append(result, alertInfo)
	}

	return result
}
