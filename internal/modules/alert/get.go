package alert

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

func (a *Alert) get(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name := luaState.Get(1)
		if name.Type() != lua.LTString {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("alert name must be a string"))
			return 2
		}

		alrt, err := a.storage.Get(name.String())
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error get alert: " + err.Error()))
			return 2
		}

		t := alrt.MarshalLua()

		luaState.Push(t)

		return 1
	}
}
