package runner

import (
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

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(j.script))
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
