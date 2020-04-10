package mock

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

type registryMock struct {
	mock.Mock
}

func (m *registryMock) Clean() {}

func (m *registryMock) Result() []modules.TestResult {
	args := m.Called()
	return args.Get(0).([]modules.TestResult)
}

func (m *registryMock) AddCall(method string, args []lua.LValue) error {
	ar := m.Called(method, args)
	return ar.Error(0)
}

func (m *registryMock) AddAssert(method string, args []lua.LValue, called bool) error {
	ar := m.Called(method, args, called)
	return ar.Error(0)
}

func (m *registryMock) Register(AnyValue, method string, callArgs, retArgs []lua.LValue) error {
	args := m.Called(AnyValue, method, callArgs, retArgs)
	return args.Error(0)
}

func (m *registryMock) Response(AnyValue, method string, args []lua.LValue) ([]lua.LValue, error) {
	ar := m.Called(AnyValue, method, args)
	return ar.Get(0).([]lua.LValue), ar.Error(1)
}

func TestOn_without_args(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := &ModuleMock{
		logger: logger,
	}

	L := lua.NewState()

	n := m.on(L)

	assert.Equal(t, 0, n)
	assert.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, len(m.errors))

	assert.Equal(t, 1, logs.FilterMessage("mock.on should have first argument").Len())
	assert.Equal(t, "mock.on should have first argument", m.errors[0])
}

func TestOn_first_arg_not_string(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := &ModuleMock{
		logger: logger,
	}

	L := lua.NewState()
	L.Push(lua.LNumber(42))

	n := m.on(L)

	assert.Equal(t, 0, n)
	assert.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, len(m.errors))

	assert.Equal(t, 1, logs.FilterMessage("mock.on first argument should be a string").Len())
	assert.Equal(t, "mock.on first argument should be a string", m.errors[0])
}

func TestOn_first_arg_is_empty(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := &ModuleMock{
		logger: logger,
	}

	L := lua.NewState()
	L.Push(lua.LString(" "))

	n := m.on(L)

	assert.Equal(t, 0, n)
	assert.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, len(m.errors))

	assert.Equal(t, 1, logs.FilterMessage("mock.on first argument should be not empty").Len())
	assert.Equal(t, "mock.on first argument should be not empty", m.errors[0])
}

func TestOn(t *testing.T) {

	// Testing: mock.on("foo", 42).response("bar", 50)

	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	reg := &registryMock{}
	reg.On("Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	m := &ModuleMock{
		logger:   logger,
		registry: reg,
	}

	L := lua.NewState()
	L.Push(lua.LString("foo"))
	L.Push(lua.LNumber(42))

	n := m.on(L)
	require.Equal(t, 1, n)

	t1 := L.Get(3).(*lua.LTable)

	L2 := lua.NewState()
	L2.Push(lua.LString("bar"))
	L2.Push(lua.LNumber(50))

	f := t1.RawGet(lua.LString("response")).(*lua.LFunction)
	n = f.GFunction(L2)

	require.Equal(t, 0, n)

	reg.AssertCalled(t, "Register", AnyValue, "foo", []lua.LValue{lua.LNumber(42)}, []lua.LValue{lua.LString("bar"), lua.LNumber(50)})
}
