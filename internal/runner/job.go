package runner

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/modules/api"
	"github.com/balerter/balerter/internal/modules/meta"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"sync/atomic"
	"time"
)

// Job represents script Job
type Job struct {
	running int64

	entryID  cron.EntryID
	name     string
	logger   *zap.Logger
	script   *script.Script
	luaState *lua.LState

	priorExecutionTime time.Duration
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

func (j *Job) GetPriorExecutionTime() time.Duration {
	return j.priorExecutionTime
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

	start := time.Now()

	err := j.luaState.DoString(string(j.script.Body))
	if err != nil {
		j.logger.Error("error run job", zap.String("script name", j.script.Name), zap.Error(err))
	}

	j.priorExecutionTime = time.Since(start)
}

func (rnr *Runner) createLuaState(j job, apiRequest *http.Request) error {
	rnr.logger.Debug("create job luaState", zap.String("name", j.Name()))

	L := lua.NewState()

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(j))
	}

	// Init storages
	for _, module := range rnr.storagesManager.Get() {
		moduleName := "storage." + module.Name()
		rnr.logger.Debug("add storage module", zap.String("name", moduleName))

		loader := module.GetLoader(j)
		L.PreloadModule(moduleName, loader)
	}

	// Init datasources
	for _, module := range rnr.dsManager.Get() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add datasource module", zap.String("name", moduleName))

		loader := module.GetLoader(j)
		L.PreloadModule(moduleName, loader)
	}

	moduleMeta := meta.New(rnr.logger)
	L.PreloadModule(moduleMeta.Name(), moduleMeta.GetLoader(j))

	a := api.New()
	err := a.FillData(apiRequest)
	if err != nil {
		return fmt.Errorf("error init api module, %w", err)
	}
	L.PreloadModule(a.Name(), a.GetLoader(j))

	j.SetLuaState(L)

	return nil
}
