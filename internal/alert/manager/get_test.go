package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestGet_NoArgs(t *testing.T) {
	m := &Manager{}

	f := m.get(nil)

	L := lua.NewState()

	n := f(L)

	assert.Equal(t, 2, n)

	e1 := L.Get(1)
	e2 := L.Get(2)

	assert.Equal(t, lua.LTNil, e1.Type())
	assert.Equal(t, "alert name must be a string", e2.String())
}

func TestGet_AlertNameNotString(t *testing.T) {
	m := &Manager{}

	f := m.get(nil)

	L := lua.NewState()
	L.Push(lua.LNumber(42))

	n := f(L)

	assert.Equal(t, 2, n)

	e1 := L.Get(2)
	e2 := L.Get(3)

	assert.Equal(t, lua.LTNil, e1.Type())
	assert.Equal(t, "alert name must be a string", e2.String())
}

func TestGet_AlertNameEmptyString(t *testing.T) {
	m := &Manager{}

	f := m.get(nil)

	L := lua.NewState()
	L.Push(lua.LString(" "))

	n := f(L)

	assert.Equal(t, 2, n)

	e1 := L.Get(2)
	e2 := L.Get(3)

	assert.Equal(t, lua.LTNil, e1.Type())
	assert.Equal(t, "alert name must be not empty", e2.String())
}

func TestGet_ErrorGet(t *testing.T) {
	mck := coreStorage.NewMock("")
	mck.AlertMock().On("Get", "foo").Return(nil, fmt.Errorf("error1"))

	m := &Manager{
		engine: mck,
	}

	f := m.get(nil)

	L := lua.NewState()
	L.Push(lua.LString("foo"))

	n := f(L)

	assert.Equal(t, 2, n)

	e1 := L.Get(2)
	e2 := L.Get(3)

	assert.Equal(t, lua.LTNil, e1.Type())
	assert.Equal(t, "error get alert: error1", e2.String())
}

func TestGet(t *testing.T) {
	a := &alert.Alert{}
	a.SetName("bar")
	a.UpdateLevel(alert.LevelError)
	a.Inc()

	mck := coreStorage.NewMock("")
	mck.AlertMock().On("Get", "foo").Return(a, nil)

	m := &Manager{
		engine: mck,
	}

	f := m.get(nil)

	L := lua.NewState()
	L.Push(lua.LString("foo"))

	n := f(L)

	assert.Equal(t, 1, n)

	e1 := L.Get(2)

	require.Equal(t, lua.LTTable, e1.Type())

	e2 := e1.(*lua.LTable)

	assert.Equal(t, "bar", e2.RawGetString("name").String())
	assert.Equal(t, "error", e2.RawGetString("level").String())
	assert.Equal(t, "1", e2.RawGetString("count").String())
}
