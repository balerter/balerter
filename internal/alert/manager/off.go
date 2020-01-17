package manager

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Manager) off(L *lua.LState) int {
	alertName, ok := m.getAlertName(L)
	if !ok {
		return 0
	}

	alertText := L.Get(2).String()

	m.activeMx.Lock()
	defer m.activeMx.Unlock()
	if _, ok := m.active[alertName]; ok {
		delete(m.active, alertName)
		m.sendSuccess(alertName, alertText)
	}

	m.logger.Debug("call alert OFF", zap.String("alertName", alertName))

	return 0
}
