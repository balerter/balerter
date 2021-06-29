package modules

import (
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
)

// TestResult represents test result
type TestResult struct {
	ScriptName string `json:"script"`
	ModuleName string `json:"module"`
	Message    string `json:"message"`
	Ok         bool   `json:"ok"`
}

// Module is an interface for core module
type Module interface {
	Name() string
	GetLoader(script *script.Script) lua.LGFunction
	Stop() error
}

// ModuleTest is an interface for core test module
type ModuleTest interface {
	Name() string
	GetLoader(script *script.Script) lua.LGFunction
	Result() ([]TestResult, error)
	Clean()
}

// ModuleMock is the module mock
type ModuleMock struct {
	mock.Mock
}

// Name returns the module name
func (m *ModuleMock) Name() string {
	args := m.Called()
	return args.String(0)
}

// GetLoader returns the lua loader
func (m *ModuleMock) GetLoader(s *script.Script) lua.LGFunction {
	args := m.Called(s)
	return args.Get(0).(lua.LGFunction)
}

// Stop the module
func (m *ModuleMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// Result returns test results
func (m *ModuleMock) Result() ([]TestResult, error) {
	args := m.Called()
	res := args.Get(0)
	if res != nil {
		return res.([]TestResult), args.Error(1)
	}
	return nil, args.Error(1)
}

// Clean the module
func (m *ModuleMock) Clean() {
	m.Called()
}
