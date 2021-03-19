package runtime

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type Runner interface {
	RunScript(name string, req *http.Request) error
}

type Runtime struct {
	runner Runner
	logger *zap.Logger
}

func New(runner Runner, logger *zap.Logger) *Runtime {
	a := &Runtime{
		runner: runner,
		logger: logger,
	}

	return a
}

func (rt *Runtime) Handler(r chi.Router) {
	r.Post("/run/{name}", rt.handlerRun)
}
