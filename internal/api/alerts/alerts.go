package alerts

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// ChManager represents interface of the Channel Manager
type ChManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

// Alerts represents alerts api module
type Alerts struct {
	alertManager corestorage.Alert
	chManager    ChManager
	logger       *zap.Logger
}

// New creates new Alerts API module
func New(alertManager corestorage.Alert, chManager ChManager, logger *zap.Logger) *Alerts {
	a := &Alerts{
		alertManager: alertManager,
		chManager:    chManager,
		logger:       logger,
	}

	return a
}

// Handler creates API handlers for Alerts API module
func (a *Alerts) Handler(r chi.Router) {
	r.Get("/", a.handlerIndex)
	r.Post("/{name}", a.handlerUpdate)
	r.Get("/{name}", a.handlerGet)
}
