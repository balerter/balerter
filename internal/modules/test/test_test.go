package test

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

type moduleTestMock struct {
	name string
	mock.Mock
}

func (m *moduleTestMock) Name() string {
	a := m.Called()
	return a.String(0)
}

func (m *moduleTestMock) GetLoader(j modules.Job) lua.LGFunction {
	a := m.Called(j)
	v := a.Get(0)
	if v == nil {
		return nil
	}
	return v.(lua.LGFunction)
}

func (m *moduleTestMock) Result() ([]modules.TestResult, error) {
	args := m.Called()
	a0 := args.Get(0)
	if a0 == nil {
		return nil, args.Error(1)
	}
	return a0.([]modules.TestResult), args.Error(1)
}

func (m *moduleTestMock) Clean() {
	m.Called()
}

type mockModulesManager struct {
	mock.Mock
}

func (m *mockModulesManager) Get() []modules.ModuleTest {
	a := m.Called()
	v := a.Get(0)
	if v == nil {
		return nil
	}
	return v.([]modules.ModuleTest)
}

func TestNew(t *testing.T) {
	dsMock := &mockModulesManager{}
	storageMock := &mockModulesManager{}

	dsModule := &moduleTestMock{}
	dsModule.On("Name").Return("m1")
	dsMock.On("Get").Return([]modules.ModuleTest{dsModule})

	storageModule := &moduleTestMock{}
	storageModule.On("Name").Return("m2")
	storageMock.On("Get").Return([]modules.ModuleTest{storageModule})

	m := New(dsMock, storageMock, nil, nil)
	assert.IsType(t, &Test{}, m)
	assert.Equal(t, 1, len(m.datasource))
	assert.Equal(t, 1, len(m.storage))

	_, ok := m.datasource["m1"]
	assert.True(t, ok)

	_, ok = m.storage["m2"]
	assert.True(t, ok)
}

func TestTest_Name(t *testing.T) {
	tst := &Test{}
	assert.Equal(t, "test", tst.Name())
}

func TestTest_Stop(t *testing.T) {
	tst := &Test{}
	assert.NoError(t, tst.Stop())
}

func TestTest_Result(t *testing.T) {
	tst := &Test{}
	v1, err := tst.Result()
	assert.Nil(t, v1)
	assert.Nil(t, err)
}

func TestTest_Clean(t *testing.T) {
	tst := &Test{}
	tst.Clean()
}

func TestTest_GetLoader(t *testing.T) {
	tst := &Test{}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := tst.GetLoader(j)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range []string{"datasource", "storage"} {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func Test_getModule_no_arg(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	tst := &Test{
		logger: zap.New(core),
	}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := tst.getModule("foo", j)

	luaState := lua.NewState()
	n := f(luaState)
	assert.Equal(t, 0, n)

	assert.Equal(t, 1, logs.FilterMessage("module should have 1 argument").Len())
}

func Test_getModule_empty_arg(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	tst := &Test{
		logger: zap.New(core),
	}
	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}
	f := tst.getModule("foo", j)

	luaState := lua.NewState()
	luaState.Push(lua.LString(" "))
	n := f(luaState)
	assert.Equal(t, 0, n)

	assert.Equal(t, 1, logs.FilterMessage("module should have 1 not empty argument").Len())
}

func Test_getModule_no_storage(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	tst := &Test{
		logger: zap.New(core),
	}
	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}
	f := tst.getModule("foo", j)

	luaState := lua.NewState()
	luaState.Push(lua.LString("bar"))
	n := f(luaState)
	assert.Equal(t, 0, n)

	assert.Equal(t, 1, logs.FilterMessage("module not found").Len())
}

func Test_getModule(t *testing.T) {
	mt := &moduleTestMock{}
	fn := func(ls *lua.LState) int {
		return 0
	}

	mt.On("GetLoader", mock.Anything).Return(lua.LGFunction(fn))

	tst := &Test{
		logger:  zap.NewNop(),
		storage: map[string]modules.ModuleTest{"bar": mt},
	}
	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}
	f := tst.getModule("test.storage", j)

	luaState := lua.NewState()
	luaState.Push(lua.LString("bar"))
	n := f(luaState)
	assert.Equal(t, 1, n)
}
