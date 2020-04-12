package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (m *Manager) get(_ *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {

		name := L.Get(1)
		if name.Type() != lua.LTString {
			L.Push(lua.LNil)
			L.Push(lua.LString("alert name must be a string"))
			return 2
		}

		nameStr := strings.TrimSpace(name.String())

		if nameStr == "" {
			L.Push(lua.LNil)
			L.Push(lua.LString("alert name must be not empty"))
			return 2
		}

		a, err := m.engine.Alert().Get(nameStr)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("error get alert: " + err.Error()))
			return 2
		}

		t := a.MarshalLua()
		alert.ReleaseAlert(a)

		L.Push(t)

		return 1
	}
}
