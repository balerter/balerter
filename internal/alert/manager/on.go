package manager

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Manager) on(L *lua.LState) int {
	alertName, ok := m.getAlertName(L)
	if !ok {
		return 0
	}

	alertText := L.Get(2).String()

	m.activeMx.Lock()
	defer m.activeMx.Unlock()
	if _, ok := m.active[alertName]; !ok {
		m.active[alertName] = 0
		m.sendError(alertName, alertText)
	}
	m.active[alertName]++

	m.logger.Debug("call alert ON", zap.String("alertName", alertName), zap.Int("count", m.active[alertName]))

	return 0
}
