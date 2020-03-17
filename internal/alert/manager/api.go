package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"time"
)

type AlertInfo struct {
	Name       string      `json:"name"`
	Level      alert.Level `json:"level"`
	LastChange time.Time   `json:"last_change"`
	Count      int         `json:"count"`
}

func (m *Manager) GetAlerts() ([]*AlertInfo, error) {

	result := make([]*AlertInfo, 0)

	alerts, err := m.engine.All()
	if err != nil {
		return nil, fmt.Errorf("error get alerts, %w", err)
	}

	for _, a := range alerts {
		alertInfo := &AlertInfo{
			Name:       a.Name(),
			Level:      a.Level(),
			LastChange: a.GetLastChangeTime(),
			Count:      a.Count(),
		}

		result = append(result, alertInfo)
	}

	return result, nil
}
