package runtime

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

// Runner is an interface for scripts runner
type Runner interface {
	RunScript(name string, req *http.Request) error
}

// Runtime represents Runtime API module
type Runtime struct {
	runner Runner
	logger *zap.Logger
}

// New creates new Runtime API module
func New(runner Runner, logger *zap.Logger) *Runtime {
	a := &Runtime{
		runner: runner,
		logger: logger,
	}

	return a
}

// Handler creates handlers for Runtime API module
func (rt *Runtime) Handler(r chi.Router) {
	r.Post("/run/{name}", rt.handlerRun)
}
