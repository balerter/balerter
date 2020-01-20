package manager

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Manager) on(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		alertName, ok := m.getAlertName(L)
		if !ok {
			return 0
		}

		alertText := L.Get(2).String()

		m.activeMx.Lock()
		defer m.activeMx.Unlock()
		info, ok := m.active[alertName]
		if !ok {
			info = &alertInfo{
				ScriptName: s.Name,
			}
			m.active[alertName] = info
			m.sendError(alertName, alertText)
		}
		info.Count++

		m.logger.Debug("call alert ON", zap.String("alertName", alertName), zap.Int("count", info.Count), zap.String("scriptName", s.Name))

		return 0
	}
}
