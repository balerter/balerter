package runner

import (
	"context"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"sync"
	"time"
)

//go:generate moq -out job_mock.go -skip-ensure -fmt goimports . job
//go:generate moq -out script_manager_mock.go -skip-ensure -fmt goimports . scriptsManager
//go:generate moq -out storages_manager_mock.go -skip-ensure -fmt goimports . storagesManager
//go:generate moq -out ds_manager_mock.go -skip-ensure -fmt goimports . dsManager

var (
	defaultUpdateInterval = time.Minute
	defaultToRunChanLen   = 64
)

type storagesManager interface {
	Get() []modules.Module
}

type scriptsManager interface {
	Get() ([]*script.Script, error)
}

type dsManager interface {
	Get() []modules.Module
}

// Runner represents the script runner
type Runner struct {
	scriptsManager  scriptsManager
	dsManager       dsManager
	storagesManager storagesManager
	cliScript       string
	logger          *zap.Logger
	updateInterval  time.Duration

	coreModules []modules.Module

	poolMx sync.Mutex
	pool   map[string]job

	cron *cron.Cron

	jobs              chan job
	updateScriptsFunc func(ctx context.Context, scripts []*script.Script, once bool)
	newJobFunc        func(s *script.Script, logger *zap.Logger) job
}

type job interface {
	Run()
	Name() string
	Script() *script.Script
	SetLuaState(ls *lua.LState)
	Stop()
	SetEntryID(cron.EntryID)
	EntryID() cron.EntryID
}

// New creates new script runner
func New(
	updateInterval time.Duration,
	scriptsManager scriptsManager,
	dsManager dsManager,
	storagesManager storagesManager,
	coreModules []modules.Module,
	cliScript string,
	logger *zap.Logger,
) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		updateInterval:  updateInterval,
		cliScript:       cliScript,
		logger:          logger,
		coreModules:     coreModules,
		pool:            make(map[string]job),
		cron:            cron.New(cron.WithSeconds(), cron.WithParser(script.CronParser)),
		jobs:            make(chan job, defaultToRunChanLen),
	}

	r.updateScriptsFunc = r.updateScripts
	r.newJobFunc = newJob

	if r.updateInterval == 0 {
		r.updateInterval = defaultUpdateInterval
	}

	go r.watchJobs()

	return r
}

func (rnr *Runner) watchJobs() {
	for j := range rnr.jobs {
		rnr.logger.Debug("run job", zap.String("name", j.Name()))
		j.Run()
	}
}

func (rnr *Runner) filterScripts(ss []*script.Script, name string) []*script.Script {
	for _, s := range ss {
		if s.Name == name {
			return []*script.Script{s}
		}
	}

	return nil
}

// Watch runs scripts watcher
func (rnr *Runner) Watch(ctx context.Context, ctxCancel context.CancelFunc, once bool) {
	rnr.cron.Start()

	defer func() {
		stopCtx := rnr.cron.Stop()
		<-stopCtx.Done()
	}()

	for {
		ss, err := rnr.scriptsManager.Get()

		// If provided CLI flag '-script', run only this script (if present)
		if rnr.cliScript != "" {
			ss = rnr.filterScripts(ss, rnr.cliScript)
		}

		if err != nil {
			rnr.logger.Error("error get scripts", zap.Error(err))
		} else {
			rnr.updateScriptsFunc(ctx, ss, once)
		}

		if once {
			ctxCancel()
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(rnr.updateInterval):
		}
	}
}

func (rnr *Runner) updateScripts(ctx context.Context, scripts []*script.Script, once bool) {
	var err error

	rnr.poolMx.Lock()
	defer rnr.poolMx.Unlock()

	rnr.logger.Debug("update scripts", zap.Int("count", len(scripts)))

	newScripts := make(map[string]struct{})

	for _, s := range scripts {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if s.Ignore {
			rnr.logger.Debug("script ignored", zap.String("name", s.Name))
			continue
		}

		newScripts[s.Hash()] = struct{}{}

		// if script already running
		if _, ok := rnr.pool[s.Hash()]; ok {
			rnr.logger.Debug("script already scheduled", zap.String("name", s.Name))
			continue
		}

		rnr.logger.Debug("schedule script job", zap.String("hash", s.Hash()), zap.String("script name", s.Name), zap.String("cron", s.CronValue))
		j := rnr.newJobFunc(s, rnr.logger)
		err = rnr.createLuaState(j, nil)
		if err != nil {
			rnr.logger.Debug("error init job", zap.String("name", s.Name), zap.Error(err))
			continue
		}

		if once {
			j.Run()
			rnr.pool[s.Hash()] = j
			continue
		}

		metrics.SetScriptsActive(j.Script().Name, true)
		f := func(j job) func() {
			return func() {
				rnr.jobs <- j
			}
		}(j)
		entryID, err := rnr.cron.AddFunc(s.CronValue, f)
		if err != nil {
			rnr.logger.Error("error schedule script", zap.String("script name", s.Name), zap.Error(err))
			continue
		}
		j.SetEntryID(entryID)

		rnr.pool[s.Hash()] = j
	}

	// stop outdated jobs
	for hash, jb := range rnr.pool {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if _, ok := newScripts[hash]; !ok {
			rnr.logger.Debug("stop script job", zap.String("hash", hash), zap.String("script name", jb.Script().Name))
			metrics.SetScriptsActive(jb.Script().Name, false)
			rnr.cron.Remove(jb.EntryID())
			jb.Stop()
			delete(rnr.pool, hash)
		}
	}
}

// Stop the module
func (rnr *Runner) Stop() {
	rnr.logger.Info("stop jobs")

	rnr.poolMx.Lock()
	defer rnr.poolMx.Unlock()

	for hash, jb := range rnr.pool {
		rnr.logger.Debug("stop script job", zap.String("hash", hash), zap.String("script name", jb.Script().Name))
		rnr.cron.Remove(jb.EntryID())
		jb.Stop()
		delete(rnr.pool, hash)
	}

	close(rnr.jobs)
}
