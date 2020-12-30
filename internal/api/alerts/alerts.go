package alerts

import (
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Alerts struct {
	alertManager corestorage.Alert
	logger       *zap.Logger
}

func New(alertManager corestorage.Alert, logger *zap.Logger) *Alerts {
	a := &Alerts{
		alertManager: alertManager,
		logger:       logger,
	}

	return a
}

func (a *Alerts) Handler(r chi.Router) {
	r.Get("/", a.handlerIndex)
	r.Post("{name}", a.handlerUpdate)
}
