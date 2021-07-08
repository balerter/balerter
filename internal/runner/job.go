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
	"sync/atomic"
)

// Job represents script Job
type Job struct {
	running int64

	entryID  cron.EntryID
	name     string
	logger   *zap.Logger
	script   *script.Script
	luaState *lua.LState
}

func newJob(s *script.Script, logger *zap.Logger) job {
	j := &Job{
		name:   s.Name,
		script: s,
		logger: logger,
	}

	return j
}

// Stop the job
func (j *Job) Stop() {
	j.luaState.Close()
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Script() *script.Script {
	return j.script
}

func (j *Job) SetLuaState(ls *lua.LState) {
	j.luaState = ls
}

func (j *Job) SetEntryID(id cron.EntryID) {
	j.entryID = id
}

func (j *Job) EntryID() cron.EntryID {
	return j.entryID
}

// Run the job
func (j *Job) Run() {
	if !atomic.CompareAndSwapInt64(&j.running, 0, 1) {
		j.logger.Debug("job already running", zap.String("name", j.name))
		return
	}
	defer atomic.StoreInt64(&j.running, 0)

	j.logger.Debug("run job", zap.String("name", j.name))

	ctx, cancel := context.WithTimeout(context.Background(), j.script.Timeout)
	defer cancel()

	j.luaState.SetContext(ctx)

	err := j.luaState.DoString(string(j.script.Body))
	if err != nil {
		j.logger.Error("error run job", zap.String("script name", j.script.Name), zap.Error(err))
	}
}

func (rnr *Runner) createLuaState(j job, apiRequest *http.Request) error {
	rnr.logger.Debug("create job luaState", zap.String("name", j.Name()))

	L := lua.NewState()

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(j.Script()))
	}

	// Init storages
	for _, module := range rnr.storagesManager.Get() {
		moduleName := "storage." + module.Name()
		rnr.logger.Debug("add storage module", zap.String("name", moduleName))

		loader := module.GetLoader(j.Script())
		L.PreloadModule(moduleName, loader)
	}

	// Init datasources
	for _, module := range rnr.dsManager.Get() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add datasource module", zap.String("name", moduleName))

		loader := module.GetLoader(j.Script())
		L.PreloadModule(moduleName, loader)
	}

	a := api.New()
	err := a.FillData(apiRequest)
	if err != nil {
		return fmt.Errorf("error init api module, %w", err)
	}
	L.PreloadModule(a.Name(), a.GetLoader(j.Script()))

	j.SetLuaState(L)

	return nil
}
