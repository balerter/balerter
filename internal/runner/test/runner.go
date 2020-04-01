package test

import (
	"github.com/balerter/balerter/internal/mock"
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
	Get() ([]*script.Script, error)
}

type Runner struct {
	scriptsManager  scriptsManager
	dsManager       manager
	storagesManager manager
	alertManager    *mock.ModuleMock
	logModule       *mock.ModuleMock
	logger          *zap.Logger

	coreModules []modules.Module
}

func New(scriptsManager scriptsManager, dsManager, storagesManager manager, alertManager, logModule *mock.ModuleMock, coreModules []modules.Module, logger *zap.Logger) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		alertManager:    alertManager,
		logModule:       logModule,
		logger:          logger,
		coreModules:     coreModules,
	}

	return r
}
