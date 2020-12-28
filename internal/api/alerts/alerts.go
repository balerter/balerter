package alerts

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Alerts struct {
	storage coreStorage.Alert
	logger  *zap.Logger
}

func New(storage coreStorage.Alert, logger *zap.Logger) *Alerts {
	a := &Alerts{
		storage: storage,
		logger:  logger,
	}

	return a
}

func (a *Alerts) Handler(r chi.Router) {
	r.Get("/", a.handlerIndex)
}
