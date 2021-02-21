package memory

import (
	"github.com/balerter/balerter/internal/alert"
)

func (m *storageAlert) Index(_ []alert.Level) (alert.Alerts, error) {
	panic("memory provider not implemented")
}

func (m *storageAlert) Update(_ string, _ alert.Level) (*alert.Alert, bool, error) {
	panic("memory provider not implemented")
}
