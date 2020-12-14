package alert

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestGet_NoArgs(t *testing.T) {
	m := &Alert{}

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
	m := &Alert{}

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

func TestGet_ErrorGet(t *testing.T) {
	err1 := fmt.Errorf("error1")

	mgrMock := &managerMock{}
	mgrMock.On("Get", mock.Anything).Return(nil, err1)

	m := &Alert{
		manager: mgrMock,
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
	a := &alert2.Alert{}
	a.SetName("bar")
	a.UpdateLevel(alert2.LevelError)
	a.Inc()

	mgrMock := &managerMock{}
	mgrMock.On("Get", mock.Anything).Return(a, nil)

	m := &Alert{
		manager: mgrMock,
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
