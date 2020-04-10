package mock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestCall(t *testing.T) {
	reg := &registryMock{}

	reg.On("AddCall", "foo", []lua.LValue{lua.LString("bar")}).Return(nil)
	reg.On("Response", AnyValue, mock.Anything, mock.Anything).Return([]lua.LValue{lua.LString("foo")}, nil)

	m := &ModuleMock{
		registry: reg,
	}

	f := m.call("foo")

	L := lua.NewState()
	L.Push(lua.LString("bar"))

	n := f(L)

	assert.Equal(t, 1, n)

	a1 := L.Get(1).(lua.LString)

	assert.Equal(t, "bar", a1.String())
}

func TestCall_Errors(t *testing.T) {
	reg := &registryMock{}

	reg.On("AddCall", "foo", []lua.LValue{lua.LString("bar")}).Return(fmt.Errorf("error1"))
	reg.On("Response", AnyValue, mock.Anything, mock.Anything).Return([]lua.LValue{lua.LString("foo")}, fmt.Errorf("error2"))

	m := &ModuleMock{
		registry: reg,
		logger:   zap.NewNop(),
	}

	f := m.call("foo")

	L := lua.NewState()
	L.Push(lua.LString("bar"))

	n := f(L)

	assert.Equal(t, 0, n)
	assert.Equal(t, 2, len(m.errors))
	assert.Equal(t, "error add query: error1", m.errors[0])
	assert.Equal(t, "error get response for method 'foo' with args '[bar]', error2", m.errors[1])
}
