package runner

import (
	"context"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"sync"
	"time"
)

type scriptsManager interface {
	Get() ([]*script.Script, error)
}

type dsManager interface {
	Get() []modules.Module
}

type alertManager interface {
	Loader() lua.LGFunction
}

type Runner struct {
	updateInterval time.Duration
	scriptsManager scriptsManager
	dsManager      dsManager
	alertManager   alertManager
	logger         *zap.Logger

	poolMx sync.Mutex
	pool   map[string]*Job
}

func New(scriptsManager scriptsManager, dsManager dsManager, alertManager alertManager, updateInterval time.Duration, logger *zap.Logger) *Runner {
	r := &Runner{
		updateInterval: updateInterval,
		scriptsManager: scriptsManager,
		dsManager:      dsManager,
		alertManager:   alertManager,
		logger:         logger,
		pool:           make(map[string]*Job),
	}

	return r
}

func (rnr *Runner) Watch(ctx context.Context, wg *sync.WaitGroup) {
	for {
		ss, err := rnr.scriptsManager.Get()

		if err != nil {
			rnr.logger.Error("error get scripts", zap.Error(err))
		} else {
			rnr.updateScripts(ctx, ss, wg)
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

		rnr.logger.Debug("run script job", zap.String("hash", s.Hash()), zap.String("script name", s.Name))
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
