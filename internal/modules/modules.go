package modules

import (
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
)

type TestResult struct {
	ScriptName string `json:"script"`
	ModuleName string `json:"module"`
	Message    string `json:"message"`
	Ok         bool   `json:"ok"`
}

type Module interface {
	Name() string
	GetLoader(script *script.Script) lua.LGFunction
	Stop() error
}

type ModuleTest interface {
	Name() string
	GetLoader(script *script.Script) lua.LGFunction
	Result() ([]TestResult, error)
	Clean()
}

type ModuleMock struct {
	mock.Mock
}

func (m *ModuleMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *ModuleMock) GetLoader(s *script.Script) lua.LGFunction {
	args := m.Called(s)
	return args.Get(0).(lua.LGFunction)
}

func (m *ModuleMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ModuleMock) Result() ([]TestResult, error) {
	args := m.Called()
	res := args.Get(0)
	if res != nil {
		return res.([]TestResult), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ModuleMock) Clean() {
	m.Called()
}
