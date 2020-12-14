package alert

import (
	lua "github.com/yuin/gopher-lua"
)

// MarshalLua marshal an Alert to the Lua table
func (a *Alert) MarshalLua() *lua.LTable {
	t := &lua.LTable{}

	t.RawSetString("name", lua.LString(a.Name()))
	t.RawSetString("level", lua.LString(a.Level().String()))
	t.RawSetString("last_change", lua.LNumber(a.GetLastChangeTime().Unix()))
	t.RawSetString("count", lua.LNumber(a.Count()))

	return t
}
