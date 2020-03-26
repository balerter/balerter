package runner

import (
	"context"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"go.uber.org/zap"
	"sync"
	"time"
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

type Runner struct {
	scriptsManager  scriptsManager
	dsManager       dsManager
	storagesManager storagesManager
	logger          *zap.Logger
	updateInterval  time.Duration

	coreModules []modules.Module

	poolMx sync.Mutex
	pool   map[string]*Job
}

func New(updateInterval time.Duration, scriptsManager scriptsManager, dsManager dsManager, storagesManager storagesManager, coreModules []modules.Module, logger *zap.Logger) *Runner {
	r := &Runner{
		scriptsManager:  scriptsManager,
		dsManager:       dsManager,
		storagesManager: storagesManager,
		updateInterval:  updateInterval,
		logger:          logger,
		coreModules:     coreModules,
		pool:            make(map[string]*Job),
	}

	if r.updateInterval == 0 {
		r.updateInterval = defaultUpdateInterval
	}

	return r
}

func (rnr *Runner) Watch(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, once bool) {
	for {
		ss, err := rnr.scriptsManager.Get()

		metrics.SetScriptsCount(len(ss))

		if err != nil {
			rnr.logger.Error("error get scripts", zap.Error(err))
		} else {
			rnr.updateScripts(ctx, ss, wg)
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

func (rnr *Runner) updateScripts(ctx context.Context, scripts []*script.Script, wg *sync.WaitGroup) {
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

		rnr.logger.Debug("run script job", zap.String("hash", s.Hash()), zap.String("script name", s.Name), zap.Duration("interval", s.Interval))
		job := newJob(s, rnr.logger)

		wg.Add(1)
		go rnr.runJob(job, wg)

		rnr.pool[s.Hash()] = job
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
			j.Stop()
			delete(rnr.pool, hash)
		}
	}
}

func (rnr *Runner) Stop() {
	rnr.logger.Debug("stop jobs")

	rnr.poolMx.Lock()
	defer rnr.poolMx.Unlock()

	for hash, j := range rnr.pool {
		rnr.logger.Debug("stop script job", zap.String("hash", hash), zap.String("script name", j.script.Name))
		j.Stop()
		delete(rnr.pool, hash)
	}
}
