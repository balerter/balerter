package test

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"go.uber.org/zap"
)

type storagesManager interface {
	Get() []modules.ModuleTest
}

type scriptsManager interface {
	Get() ([]*script.Script, error)
}

type dsManager interface {
	Get() []modules.ModuleTest
	Result() []modules.TestResult
}

type Runner struct {
	scriptsManager  scriptsManager
	dsManager       dsManager
	storagesManager storagesManager
	logger          *zap.Logger

	coreModules []modules.Module
}

func New(scriptsManager scriptsManager, dsManager dsManager, storagesManager storagesManager, coreModules []modules.Module, logger *zap.Logger) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		logger:          logger,
		coreModules:     coreModules,
	}

	return r
}
