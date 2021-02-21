package alerts

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type ChManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

type Alerts struct {
	alertManager corestorage.Alert
	chManager    ChManager
	logger       *zap.Logger
}

func New(alertManager corestorage.Alert, chManager ChManager, logger *zap.Logger) *Alerts {
	a := &Alerts{
		alertManager: alertManager,
		chManager:    chManager,
		logger:       logger,
	}

	return a
}

func (a *Alerts) Handler(r chi.Router) {
	r.Get("/", a.handlerIndex)
	r.Post("/{name}", a.handlerUpdate)
}
