package alert

import (
	alerts "github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

func (a *Alert) get(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name := luaState.Get(1) //nolint:mnd
		if name.Type() != lua.LTString {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("alert name must be a string"))
			return 2 //nolint:mnd
		}

		alrt, err := a.manager.Get(name.String())
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error get alert: " + err.Error()))
			return 2 //nolint:mnd
		}

		t := alrt.MarshalLua()
		alerts.ReleaseAlert(alrt)

		luaState.Push(t)

		return 1
	}
}
