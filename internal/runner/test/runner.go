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

type Runner struct {
	scriptsManager  scriptsManager
	dsManager       manager
	storagesManager manager
	testModule      modules.Module
	logger          *zap.Logger

	coreModules []modules.ModuleTest
}

func New(scriptsManager scriptsManager,
	dsManager,
	storagesManager manager,
	testModule modules.Module,
	coreModules []modules.ModuleTest,
	logger *zap.Logger,
) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		testModule:      testModule,
		logger:          logger,
		coreModules:     coreModules,
	}

	return r
}
