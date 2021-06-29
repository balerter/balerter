package converter

import (
	lua "github.com/yuin/gopher-lua"
)

// FromDateBytes returns lua string from bytes
func FromDateBytes(v interface{}) lua.LValue {
	t := v.(*[]byte)
	return lua.LString(*t)
}
