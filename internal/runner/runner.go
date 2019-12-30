package runner

import (
	"github.com/balerter/balerter/internal/script/script"
	"go.uber.org/zap"
)

type scriptsProvider interface {
	Get() ([]*script.Script, error)
}

type Runner struct {
	scriptsProvider scriptsProvider
	logger          *zap.Logger
}

func New(scriptsProvider scriptsProvider, logger *zap.Logger) *Runner {
	r := &Runner{
		scriptsProvider: scriptsProvider,
		logger:          logger,
	}

	return r
}

func (rnr *Runner) Run() {

}
