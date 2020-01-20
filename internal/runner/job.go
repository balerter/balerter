package runner

import (
	moduleLog "github.com/balerter/balerter/internal/modules/log"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Job struct {
	name   string
	logger *zap.Logger
	script *script.Script
	stop   chan struct{}
}

func newJob(s *script.Script, logger *zap.Logger) *Job {
	j := &Job{
		name:   s.Name,
		script: s,
		stop:   make(chan struct{}),
		logger: logger,
	}

	return j
}

func (j *Job) Stop() {
	close(j.stop)
}

func (rnr *Runner) runJob(j *Job, wg *sync.WaitGroup) {
	defer wg.Done()
	rnr.logger.Debug("run job loop", zap.String("name", j.name))

	L := rnr.createLuaState(j)
	defer L.Close()

	for {
		rnr.logger.Debug("run job", zap.String("name", j.name))
		err := L.DoString(string(j.script.Body))
		if err != nil {
			j.logger.Error("error run job", zap.String("script name", j.script.Name), zap.Error(err))
		}

		select {
		case <-j.stop:
			return
		case <-time.After(j.script.Interval):
		}
	}
}

func (rnr *Runner) createLuaState(j *Job) *lua.LState {
	rnr.logger.Debug("create job", zap.String("name", j.name))

	L := lua.NewState()

	// Init core modules
	L.PreloadModule("log", moduleLog.New(j.name, rnr.logger))
	L.PreloadModule("alert", rnr.alertManager.Loader(j.script))

	// Init datasources
	for _, module := range rnr.dsManager.Get() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add module", zap.String("name", moduleName))

		loader := module.GetLoader()
		L.PreloadModule(moduleName, loader)
	}

	return L
}
