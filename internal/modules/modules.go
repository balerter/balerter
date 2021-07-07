package modules

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

//go:generate moq -out module_mock.go -skip-ensure -fmt goimports . Module
//go:generate moq -out module_test_mock.go -skip-ensure -fmt goimports . ModuleTest

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
