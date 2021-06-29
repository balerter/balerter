package test

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"go.uber.org/zap"
)

type manager interface {
	Get() []modules.ModuleTest
	Result() ([]modules.TestResult, error)
	Clean()
}

type scriptsManager interface {
	GetWithTests() ([]*script.Script, error)
}

// Runner represents the Test runner
type Runner struct {
	scriptsManager  scriptsManager
	dsManager       manager
	storagesManager manager
	logger          *zap.Logger

	coreModules []modules.ModuleTest
}

// New creates new test runner
func New(scriptsManager scriptsManager,
	dsManager,
	storagesManager manager,
	coreModules []modules.ModuleTest,
	logger *zap.Logger,
) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		logger:          logger,
		coreModules:     coreModules,
	}

	return r
}
