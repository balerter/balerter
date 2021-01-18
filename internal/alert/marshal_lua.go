package alert

import (
	lua "github.com/yuin/gopher-lua"
)

// MarshalLua marshal an Alert to the Lua table
func (a *Alert) MarshalLua() *lua.LTable {
	t := &lua.LTable{}

	t.RawSetString("Name", lua.LString(a.Name))
	t.RawSetString("Level", lua.LString(a.Level.String()))
	t.RawSetString("last_change", lua.LNumber(a.LastChange.Unix()))
	t.RawSetString("Count", lua.LNumber(a.Count))

	return t
}
