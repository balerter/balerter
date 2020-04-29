package runner

import (
	"context"
	"sync"
	"time"

	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
)

var (
	defaultUpdateInterval = time.Minute
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

type runningJob struct {
	*Job
	ID cron.EntryID
}

type Runner struct {
	scriptsManager  scriptsManager
	dsManager       dsManager
	storagesManager storagesManager
	logger          *zap.Logger
	updateInterval  time.Duration
	cron            *cron.Cron

	coreModules []modules.Module

	poolMx sync.Mutex
	pool   map[string]*runningJob
}

func New(updateInterval time.Duration, scriptsManager scriptsManager, dsManager dsManager, storagesManager storagesManager, coreModules []modules.Module, logger *zap.Logger) *Runner {
	c := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(
				cron.DefaultLogger,
			),
		),
	)

	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		updateInterval:  updateInterval,
		cron:            c,
		logger:          logger,
		coreModules:     coreModules,
		pool:            make(map[string]*runningJob),
	}

	if r.updateInterval == 0 {
		r.updateInterval = defaultUpdateInterval
	}

	return r
}

func (rnr *Runner) Watch(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, once bool) {
	wg.Add(1)
	rnr.cron.Start()
	defer func() {
		rnr.stop()
		wg.Done()
	}()

	for {
		ss, err := rnr.scriptsManager.Get()

		metrics.SetScriptsCount(len(ss))

		if err != nil {
			rnr.logger.Error("error get scripts", zap.Error(err))
		} else {
			rnr.updateScripts(ctx, ss)
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

func (rnr *Runner) updateScripts(ctx context.Context, scripts []*script.Script) {
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

		newScripts[s.Hash()] = struct{}{}

		// if script already running
		if _, ok := rnr.pool[s.Hash()]; ok {
			rnr.logger.Debug("script already running", zap.String("name", s.Name))
			continue
		}

		rnr.logger.Debug("run script job", zap.String("hash", s.Hash()), zap.String("script name", s.Name), zap.String("schedule", s.ScheduleString))
		job := newJob(s, rnr.logger, rnr.createLuaState)
		id := rnr.cron.Schedule(s.Schedule, job)
		rnr.pool[s.Hash()] = &runningJob{
			Job: job,
			ID:  id,
		}
	}

	// stop outdated jobs
	for hash, j := range rnr.pool {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if _, ok := newScripts[hash]; !ok {
			rnr.logger.Debug("stop script job", zap.String("hash", hash), zap.String("script name", j.script.Name))
			rnr.cron.Remove(j.ID)
			delete(rnr.pool, hash)
		}
	}
}

func (rnr *Runner) stop() {
	rnr.logger.Debug("stop jobs")

	rnr.poolMx.Lock()
	defer rnr.poolMx.Unlock()

	for hash, j := range rnr.pool {
		rnr.logger.Debug("stop script job", zap.String("hash", hash), zap.String("script name", j.script.Name))
		rnr.cron.Remove(j.ID)
		delete(rnr.pool, hash)
	}

	ctx := rnr.cron.Stop()
	<-ctx.Done()
}
