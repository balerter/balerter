package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (m *Manager) get(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name := luaState.Get(1) //nolint:mnd
		if name.Type() != lua.LTString {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("alert name must be a string"))
			return 2 //nolint:mnd
		}

		nameStr := strings.TrimSpace(name.String())

		if nameStr == "" {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("alert name must be not empty"))
			return 2 //nolint:mnd
		}

		a, err := m.engine.Alert().Get(nameStr)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error get alert: " + err.Error()))
			return 2 //nolint:mnd
		}

		t := a.MarshalLua()
		alert.ReleaseAlert(a)

		luaState.Push(t)

		return 1
	}
}
