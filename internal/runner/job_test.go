package runner

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

type alertManagerMock struct {
	mock.Mock
}

func (m *alertManagerMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *alertManagerMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *alertManagerMock) GetLoader(s *script.Script) lua.LGFunction {
	args := m.Called(s)
	return args.Get(0).(lua.LGFunction)
}

type dsManagerMock struct {
	mock.Mock
}

func (m *dsManagerMock) Errors() []error {
	args := m.Called()
	return args.Get(0).([]error)
}

func (m *dsManagerMock) GetMocks() []modules.ModuleTest {
	args := m.Called()
	return args.Get(0).([]modules.ModuleTest)
}

func (m *dsManagerMock) Get() []modules.Module {
	args := m.Called()
	return args.Get(0).([]modules.Module)
}

type storagesManagerMock struct {
	mock.Mock
}

func (m *storagesManagerMock) Get() []modules.Module {
	args := m.Called()
	return args.Get(0).([]modules.Module)
}

type moduleMock struct {
	mock.Mock
}

func (m *moduleMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *moduleMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *moduleMock) GetLoader(s *script.Script) lua.LGFunction {
	args := m.Called(s)
	return args.Get(0).(lua.LGFunction)
}

func TestRunner_createLuaState(t *testing.T) {
	m1 := &moduleMock{}
	m1.On("Name").Return("module1")
	m1.On("GetLoader", mock.Anything).Return(func() lua.LGFunction {
		return func(state *lua.LState) int {
			return 0
		}
	}())

	dsManager := &dsManagerMock{}
	dsManager.On("Get").Return([]modules.Module{m1})

	storagesManager := &storagesManagerMock{}
	storagesManager.On("Get").Return([]modules.Module{m1})

	alertManager := &alertManagerMock{}
	alertManager.On("Name").Return("alert")
	alertManager.On("GetLoader", mock.Anything).Return(func() lua.LGFunction {
		return func(state *lua.LState) int {
			return 0
		}
	}())

	rnr := &Runner{
		logger:          zap.NewNop(),
		dsManager:       dsManager,
		storagesManager: storagesManager,
		coreModules:     []modules.Module{alertManager},
	}

	j := &Job{name: "job1"}

	err := rnr.createLuaState(j, nil)
	assert.NoError(t, err)

	m1.AssertCalled(t, "Name")
	m1.AssertCalled(t, "GetLoader", mock.Anything)

	require.NotNil(t, j.luaState)
	m1.AssertExpectations(t)
	dsManager.AssertExpectations(t)
}
