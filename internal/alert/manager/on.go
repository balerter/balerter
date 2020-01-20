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
	info, ok := m.active[alertName]
	if !ok {
		info = &alertInfo{}
		m.active[alertName] = info
		m.sendError(alertName, alertText)
	}
	info.count++

	m.logger.Debug("call alert ON", zap.String("alertName", alertName), zap.Int("count", info.count))

	return 0
}
