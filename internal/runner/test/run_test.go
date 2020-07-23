package test

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func Test_splitScripts(t *testing.T) {
	src := []*script.Script{
		{Name: "s1", IsTest: false},
		{Name: "s1_test", IsTest: true, TestTarget: "s1"},
		{Name: "s2", IsTest: false},
		{Name: "s2_test", IsTest: true, TestTarget: "s2"},
	}

	pairs, err := splitScripts(src)
	require.NoError(t, err)

	assert.Equal(t, 2, len(pairs))

	p, ok := pairs["s1_test"]
	require.True(t, ok)
	assert.Equal(t, "s1_test", p.test.Name)
	assert.Equal(t, "s1", p.main.Name)

	p, ok = pairs["s2_test"]
	require.True(t, ok)
	assert.Equal(t, "s2_test", p.test.Name)
	assert.Equal(t, "s2", p.main.Name)
}

func Test_splitScripts_no_main(t *testing.T) {
	src := []*script.Script{
		{Name: "s1_test", IsTest: true, TestTarget: "s1"},
	}

	_, err := splitScripts(src)
	require.Error(t, err)

	assert.Equal(t, "main script for test 's1_test' not found", err.Error())
}

type scriptManagerMock struct {
	mock.Mock
}

func (s *scriptManagerMock) GetWithTests() ([]*script.Script, error) {
	args := s.Called()
	a1 := args.Get(0)
	if a1 == nil {
		return nil, args.Error(1)
	}
	return a1.([]*script.Script), args.Error(1)
}

func Test_Run_error_get_scripts(t *testing.T) {
	m := &scriptManagerMock{}
	m.On("GetWithTests").Return(nil, fmt.Errorf("error1"))

	rnr := &Runner{
		scriptsManager: m,
	}

	_, testFail, err := rnr.Run()
	require.Error(t, err)
	assert.False(t, testFail)
	assert.Equal(t, "error get scripts, error1", err.Error())
}

func Test_Run_error_split_scripts(t *testing.T) {
	m := &scriptManagerMock{}
	m.On("GetWithTests").Return([]*script.Script{{Name: "s1", IsTest: true}}, nil)

	rnr := &Runner{
		scriptsManager: m,
	}

	_, testFail, err := rnr.Run()
	require.Error(t, err)
	assert.False(t, testFail)
	assert.Equal(t, "error select tests, main script for test 's1' not found", err.Error())
}

type managerMock struct {
	mock.Mock
}

func (m *managerMock) Get() []modules.ModuleTest {
	args := m.Called()
	return args.Get(0).([]modules.ModuleTest)
}

func (m *managerMock) Result() ([]modules.TestResult, error) {
	args := m.Called()
	return args.Get(0).([]modules.TestResult), args.Error(1)
}

func (m *managerMock) Clean() {

}

func Test_runPair(t *testing.T) {
	storageManagerMock := &managerMock{}
	storageManagerMock.On("Get").Return([]modules.ModuleTest{})
	storageManagerMock.On("Result").Return([]modules.TestResult{}, nil)

	dsManagerMock := &managerMock{}
	dsManagerMock.On("Get").Return([]modules.ModuleTest{})
	dsManagerMock.On("Result").Return([]modules.TestResult{}, nil)

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
	}

	p := pair{
		main: &script.Script{
			Name:   "s1",
			Body:   []byte("a = 10"),
			IsTest: false,
		},
		test: &script.Script{
			Name:       "s1_test",
			Body:       []byte("a = 10"),
			IsTest:     true,
			TestTarget: "s1",
		},
	}

	var res []modules.TestResult

	res, err := rnr.runPair(res, "pair", p)
	require.NoError(t, err)

	assert.Equal(t, 1, len(res))

	r := res[0]

	assert.Equal(t, true, r.Ok)
	assert.Equal(t, "PASS", r.Message)
}
