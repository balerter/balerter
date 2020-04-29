package runner

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type LuaStateGenerateFn func(*Job) *lua.LState

type Job struct {
	name          string
	logger        *zap.Logger
	script        *script.Script
	generateState LuaStateGenerateFn
}

func newJob(s *script.Script, logger *zap.Logger, stateGen LuaStateGenerateFn) *Job {
	j := &Job{
		name:          s.Name,
		script:        s,
		logger:        logger,
		generateState: stateGen,
	}

	return j
}

func (j *Job) Run() {
	L := j.generateState(j)
	defer L.Close()

	j.logger.Debug("run job", zap.String("name", j.name))
	err := L.DoString(string(j.script.Body))
	if err != nil {
		j.logger.Error("error run job", zap.String("script name", j.script.Name), zap.Error(err))
	}
}

func (rnr *Runner) createLuaState(j *Job) *lua.LState {
	rnr.logger.Debug("create job", zap.String("name", j.name))

	L := lua.NewState()

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(j.script))
	}

	// Init storages
	for _, module := range rnr.storagesManager.Get() {
		moduleName := "storage." + module.Name()
		rnr.logger.Debug("add storage module", zap.String("name", moduleName))

		loader := module.GetLoader(j.script)
		L.PreloadModule(moduleName, loader)
	}

	// Init datasources
	for _, module := range rnr.dsManager.Get() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add datasource module", zap.String("name", moduleName))

		loader := module.GetLoader(j.script)
		L.PreloadModule(moduleName, loader)
	}

	return L
}
