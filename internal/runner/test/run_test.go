package test

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
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

func Test_runPair(t *testing.T) {
	mTest1 := &modules.ModuleTestMock{
		NameFunc: func() string {
			return "mt1"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return getLGFunc()
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {},
	}

	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{mTest1}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{{Ok: true, Message: "PASS"}}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{mTest1}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{{Ok: true, Message: "PASS"}}, nil
		},
		CleanFunc: func() {

		},
	}

	coreModule1 := &modules.ModuleTestMock{
		NameFunc: func() string {
			return "m1"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return getLGFunc()
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{{Ok: true, Message: "PASS"}}, nil
		},
		CleanFunc: func() {

		},
	}

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
		coreModules:     []modules.ModuleTest{coreModule1},
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

	assert.Equal(t, 4, len(res))

	for _, r := range res {
		assert.Equal(t, true, r.Ok)
		assert.Equal(t, "PASS", r.Message)
	}
}

func Test_runPair_with_fail_result(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{{Ok: false, Message: "FAIL"}}, nil
		},
		CleanFunc: func() {

		},
	}

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

	assert.Equal(t, 2, len(res))

	r := res[0]
	assert.Equal(t, false, r.Ok)
	assert.Equal(t, "FAIL", r.Message)

	r = res[1]
	assert.Equal(t, false, r.Ok)
	assert.Equal(t, "FAIL", r.Message)

	assert.Equal(t, 1, len(dsManagerMock.CleanCalls()))
	assert.Equal(t, 1, len(dsManagerMock.ResultCalls()))
	assert.Equal(t, 2, len(dsManagerMock.GetCalls()))
}

func Test_runPair_error_run_test(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
	}

	p := pair{
		test: &script.Script{
			Name:       "s1_test",
			Body:       []byte{0x00},
			IsTest:     true,
			TestTarget: "s1",
		},
	}

	_, err := rnr.runPair([]modules.TestResult{}, "pair", p)
	require.Error(t, err)
	assert.Equal(t, "error run test job, <string> line:1(column:1) near '\x00':   Invalid token\n", err.Error())
}

func Test_runPair_error_run_main(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
	}

	p := pair{
		main: &script.Script{
			Name:   "s1",
			Body:   []byte{0x00},
			IsTest: false,
		},
		test: &script.Script{
			Name:       "s1_test",
			Body:       []byte("a = 10"),
			IsTest:     true,
			TestTarget: "s1",
		},
	}

	_, err := rnr.runPair([]modules.TestResult{}, "pair", p)
	require.Error(t, err)
	assert.Equal(t, "error run main job, <string> line:1(column:1) near '\x00':   Invalid token\n", err.Error())
}

func Test_runPair_error_get_ds_results(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, fmt.Errorf("error1")
		},
		CleanFunc: func() {

		},
	}

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

	_, err := rnr.runPair([]modules.TestResult{}, "pair", p)
	require.Error(t, err)
	assert.Equal(t, "error get results from datasource manager, error1", err.Error())
}

func Test_runPair_error_get_storages_results(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, fmt.Errorf("error1")
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

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

	_, err := rnr.runPair([]modules.TestResult{}, "pair", p)
	require.Error(t, err)
	assert.Equal(t, "error get results from storage manager, error1", err.Error())
}

func getLGFunc() lua.LGFunction {
	return func(_ *lua.LState) int {
		return 0
	}
}

func Test_runPair_error_get_core_module_results(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	coreModule1 := &modules.ModuleTestMock{
		NameFunc: func() string {
			return "m1"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return getLGFunc()
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, fmt.Errorf("error1")
		},
	}

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
		coreModules: []modules.ModuleTest{
			coreModule1,
		},
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

	_, err := rnr.runPair([]modules.TestResult{}, "pair", p)
	require.Error(t, err)
	assert.Equal(t, "error get results from 'm1' module, error1", err.Error())
}

func Test_Run(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	scriptMgrMock := &scriptManagerMock{}
	scriptMgrMock.On("GetWithTests").Return([]*script.Script{
		{Name: "s1", Body: []byte("a = 10"), IsTest: false},
		{Name: "s1_test", Body: []byte("a = 10"), IsTest: true, TestTarget: "s1"},
	}, nil)

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
		scriptsManager:  scriptMgrMock,
	}

	res, ok, err := rnr.Run()
	require.NoError(t, err)

	assert.True(t, ok)
	assert.Equal(t, 1, len(res))
}

func Test_Run_error_runPair(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	scriptMgrMock := &scriptManagerMock{}
	scriptMgrMock.On("GetWithTests").Return([]*script.Script{
		{Name: "s1", Body: []byte{0x00}, IsTest: false},
		{Name: "s1_test", Body: []byte("a = 10"), IsTest: true, TestTarget: "s1"},
	}, nil)

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
		scriptsManager:  scriptMgrMock,
	}

	_, ok, err := rnr.Run()
	require.Error(t, err)
	assert.Equal(t, "error run main job, <string> line:1(column:1) near '\x00':   Invalid token\n", err.Error())
	assert.False(t, ok)
}

func Test_Run_fail_tests(t *testing.T) {
	storageManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{{Ok: false, Message: "FAIL"}}, nil
		},
		CleanFunc: func() {

		},
	}

	dsManagerMock := &managerMock{
		GetFunc: func() []modules.ModuleTest {
			return []modules.ModuleTest{}
		},
		ResultFunc: func() ([]modules.TestResult, error) {
			return []modules.TestResult{}, nil
		},
		CleanFunc: func() {

		},
	}

	scriptMgrMock := &scriptManagerMock{}
	scriptMgrMock.On("GetWithTests").Return([]*script.Script{
		{Name: "s1", Body: []byte("a = 10"), IsTest: false},
		{Name: "s1_test", Body: []byte("a = 10"), IsTest: true, TestTarget: "s1"},
	}, nil)

	rnr := &Runner{
		logger:          zap.NewNop(),
		storagesManager: storageManagerMock,
		dsManager:       dsManagerMock,
		scriptsManager:  scriptMgrMock,
	}

	res, ok, err := rnr.Run()
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 2, len(res))
}
