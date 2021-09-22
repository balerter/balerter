package modules

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"time"
)

//go:generate moq -out module_mock.go -skip-ensure -fmt goimports . Module
//go:generate moq -out module_test_mock.go -skip-ensure -fmt goimports . ModuleTest
//go:generate moq -out job_mock.go -skip-ensure -fmt goimports . Job

// TestResult represents test result
type TestResult struct {
	ScriptName string `json:"script"`
	ModuleName string `json:"module"`
	Message    string `json:"message"`
	Ok         bool   `json:"ok"`
}

type Job interface {
	Script() *script.Script
	GetPriorExecutionTime() time.Duration
}

// Module is an interface for core module
type Module interface {
	Name() string
	GetLoader(j Job) lua.LGFunction
	Stop() error
}

// ModuleTest is an interface for core test module
type ModuleTest interface {
	Name() string
	GetLoader(j Job) lua.LGFunction
	Result() ([]TestResult, error)
	Clean()
}
