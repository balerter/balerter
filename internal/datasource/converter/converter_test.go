package converter

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestConvertFunctions(t *testing.T) {
	v := sql.NullFloat64{Float64: 1.12, Valid: true}
	assert.Equal(t, lua.LNumber(1.12), FromFloat64(&v))
	v.Valid = false
	assert.Equal(t, lua.LNil, FromFloat64(&v))

	d := sql.NullTime{Time: time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC), Valid: true}
	assert.Equal(t, lua.LString("2020-02-03"), FromDate(&d))
	d.Valid = false
	assert.Equal(t, lua.LNil, FromDate(&d))

	d = sql.NullTime{Time: time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC), Valid: true}
	assert.Equal(t, lua.LString("1580702706"), FromDateTime(&d))
	d.Valid = false
	assert.Equal(t, lua.LNil, FromDateTime(&d))

	s := sql.NullString{String: "value", Valid: true}
	assert.Equal(t, lua.LString("value"), FromString(&s))
	s.Valid = false
	assert.Equal(t, lua.LNil, FromString(&s))

	b := sql.NullBool{Bool: true, Valid: true}
	assert.Equal(t, lua.LBool(true), FromBoolean(&b))
	b.Valid = false
	assert.Equal(t, lua.LNil, FromBoolean(&b))

	b = sql.NullBool{Bool: false, Valid: true}
	assert.Equal(t, lua.LBool(false), FromBoolean(&b))
	b.Valid = false
	assert.Equal(t, lua.LNil, FromBoolean(&b))

	ui := uint(100)
	assert.Equal(t, lua.LNumber(100), FromUInt(&ui))

	i := sql.NullInt64{Int64: 100, Valid: true}
	assert.Equal(t, lua.LNumber(100), FromInt(&i))
	i.Valid = false
	assert.Equal(t, lua.LNil, FromInt(&i))

	i = sql.NullInt64{Int64: -100, Valid: true}
	assert.Equal(t, lua.LNumber(-100), FromInt(&i))
	i.Valid = false
	assert.Equal(t, lua.LNil, FromInt(&i))
}
