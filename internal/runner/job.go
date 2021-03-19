package runner

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/modules/api"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
)

type Job struct {
	entryID  cron.EntryID
	name     string
	logger   *zap.Logger
	script   *script.Script
	luaState *lua.LState
}

func newJob(s *script.Script, logger *zap.Logger) *Job {
	j := &Job{
		name:   s.Name,
		script: s,
		logger: logger,
	}

	return j
}

func (j *Job) Stop() {
	j.luaState.Close()
}

func (j *Job) Run() {
	j.logger.Debug("run job", zap.String("name", j.name))

	ctx, cancel := context.WithTimeout(context.Background(), j.script.Timeout)
	defer cancel()

	j.luaState.SetContext(ctx)

	err := j.luaState.DoString(string(j.script.Body))
	if err != nil {
		j.logger.Error("error run job", zap.String("script name", j.script.Name), zap.Error(err))
	}
}

func (rnr *Runner) createLuaState(j *Job, apiRequest *http.Request) error {
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

	a := api.New()
	err := a.FillData(apiRequest)
	if err != nil {
		return fmt.Errorf("error init api module, %w", err)
	}
	L.PreloadModule(a.Name(), a.GetLoader(j.script))

	j.luaState = L

	return nil
}
