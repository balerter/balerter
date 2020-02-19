package converter

import (
	lua "github.com/yuin/gopher-lua"
	"time"
)

func FromFloat64(v interface{}) lua.LValue {
	return lua.LNumber(*(v.(*float64)))
}

func FromDate(v interface{}) lua.LValue {
	return lua.LString(v.(*time.Time).Format("2006-01-02"))
}

func FromDateTime(v interface{}) lua.LValue {
	return lua.LString(v.(*time.Time).Format(time.RFC3339))
}

func FromString(v interface{}) lua.LValue {
	return lua.LString(*(v.(*string)))
}

func FromBoolean(v interface{}) lua.LValue {
	return lua.LBool(*v.(*bool))
}

func FromUInt(v interface{}) lua.LValue {
	return lua.LNumber(int(*(v.(*uint))))
}

func FromInt(v interface{}) lua.LValue {
	return lua.LNumber(int(*(v.(*int))))
}

func FromTimestamp(v interface{}) lua.LValue {
	return lua.LString("<NOT IMPLEMENTED>")
}
