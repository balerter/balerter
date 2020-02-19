package converter

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestConvertFunctions(t *testing.T) {
	v := 1.12
	assert.Equal(t, lua.LNumber(1.12), FromFloat64(&v))

	d := time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC)
	assert.Equal(t, lua.LString("2020-02-03"), FromDate(&d))

	d = time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC)
	assert.Equal(t, lua.LString("2020-02-03T04:05:06Z"), FromDateTime(&d))

	s := "value"
	assert.Equal(t, lua.LString("value"), FromString(&s))

	b := true
	assert.Equal(t, lua.LBool(true), FromBoolean(&b))

	b = false
	assert.Equal(t, lua.LBool(false), FromBoolean(&b))

	ui := uint(100)
	assert.Equal(t, lua.LNumber(100), FromUInt(&ui))

	i := 100
	assert.Equal(t, lua.LNumber(100), FromInt(&i))

	i = -100
	assert.Equal(t, lua.LNumber(-100), FromInt(&i))
}
