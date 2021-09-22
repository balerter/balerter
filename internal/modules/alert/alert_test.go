package alert

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "alert", ModuleName())
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{
		"warn",
		"warning",
		"error",
		"fail",
		"success",
		"ok",
		"get",
	}, Methods())
}

func TestNew(t *testing.T) {
	a := New(nil, nil, nil)
	assert.IsType(t, &Alert{}, a)
}

func TestName(t *testing.T) {
	a := &Alert{}
	assert.Equal(t, "alert", a.Name())
}

func TestAlert_GetLoader(t *testing.T) {
	a := &Alert{}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := a.GetLoader(j)

	ls := lua.NewState()
	n := f(ls)
	assert.Equal(t, 1, n)

	v := ls.Get(1).(*lua.LTable)

	for _, m := range Methods() {
		assert.Equal(t, lua.LTFunction, v.RawGetString(m).Type())
	}
}

func TestStop(t *testing.T) {
	a := &Alert{}
	assert.NoError(t, a.Stop())
}
