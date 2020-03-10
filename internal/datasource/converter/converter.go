package converter

import (
	lua "github.com/yuin/gopher-lua"
)

func FromDateBytes(v interface{}) lua.LValue {
	t := v.(*[]byte)
	return lua.LString(string(*t))
}
