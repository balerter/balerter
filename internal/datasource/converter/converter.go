package converter

import (
	"database/sql"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

func FromFloat64(v interface{}) lua.LValue {
	vt := v.(*sql.NullFloat64)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LNumber(vt.Float64)
}

func FromDate(v interface{}) lua.LValue {
	vt := v.(*sql.NullTime)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LString(vt.Time.Format("2006-01-02"))
}

func FromDateTime(v interface{}) lua.LValue {
	vt := v.(*sql.NullTime)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LString(strconv.Itoa(int(vt.Time.Unix())))
}

func FromString(v interface{}) lua.LValue {
	vt := v.(*sql.NullString)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LString(vt.String)
}

func FromBoolean(v interface{}) lua.LValue {
	vt := v.(*sql.NullBool)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LBool(vt.Bool)
}

func FromUInt(v interface{}) lua.LValue {
	// todo Use sql.NullInt64 ?
	return lua.LNumber(int(*(v.(*uint))))
}

func FromInt(v interface{}) lua.LValue {
	vt := v.(*sql.NullInt64)
	if !vt.Valid {
		return lua.LNil
	}
	return lua.LNumber(int(vt.Int64))
}
