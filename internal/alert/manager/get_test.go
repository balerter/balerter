package manager

import (
	"github.com/stretchr/testify/assert"
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
