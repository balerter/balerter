package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

type moduleTestMock struct {
	name string
	mock.Mock
}

func (m *moduleTestMock) Name() string {
	m.Called()
	return ""
}

func (m *moduleTestMock) GetLoader(_ modules.Job) lua.LGFunction {
	m.Called()
	return nil
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

func TestNew(t *testing.T) {
	m := New(zap.NewNop())
	assert.IsType(t, &Manager{}, m)
}

func TestManager_Get(t *testing.T) {
	m1 := &moduleTestMock{name: "foo"}
	m2 := &moduleTestMock{name: "bar"}

	m := &Manager{
		logger:  zap.NewNop(),
		modules: map[string]modules.ModuleTest{"foo": m1, "bar": m2},
	}

	result := m.Get()
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, m1)
	assert.Contains(t, result, m2)
}

func TestManager_Init(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.ModuleTest{},
	}

	err := m.Init(&upload.Upload{S3: []s3.S3{{
		Name:     "f1",
		Region:   "f2",
		Key:      "f3",
		Secret:   "f4",
		Endpoint: "f5",
		Bucket:   "f6",
	}}})

	require.NoError(t, err)

	md, ok := m.modules["s3.f1"]
	require.True(t, ok)

	assert.Equal(t, "s3.f1", md.Name())
}

func TestManager_Clean(t *testing.T) {
	m1 := &moduleTestMock{name: "foo"}
	m2 := &moduleTestMock{name: "bar"}

	m1.On("Clean").Return()
	m2.On("Clean").Return()

	m := &Manager{
		logger:  zap.NewNop(),
		modules: map[string]modules.ModuleTest{"foo": m1, "bar": m2},
	}

	m.Clean()

	m1.AssertCalled(t, "Clean")
	m2.AssertCalled(t, "Clean")
}

func TestManager_Result_error_from_module(t *testing.T) {
	m1 := &moduleTestMock{name: "foo"}
	m2 := &moduleTestMock{name: "bar"}

	e1 := fmt.Errorf("err1")

	m1.On("Result").Return(nil, nil)
	m2.On("Result").Return(nil, e1)

	m := &Manager{
		logger:  zap.NewNop(),
		modules: map[string]modules.ModuleTest{"foo": m1, "bar": m2},
	}

	_, err := m.Result()
	assert.Equal(t, e1, err)
}
