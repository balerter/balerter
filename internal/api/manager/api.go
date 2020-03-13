package manager

import (
	"context"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/balerter/balerter/internal/api/alerts"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

type alertManagerAPIer interface {
	GetAlerts() []*alertManager.AlertInfo
}

type API struct {
	address string
	server  *http.Server
	logger  *zap.Logger
}

func New(cfg *config.API, alertManager alertManagerAPIer, logger *zap.Logger) *API {
	api := &API{
		address: cfg.Address,
		server:  &http.Server{},
		logger:  logger,
	}

	m := http.NewServeMux()

	m.HandleFunc("/api/v1/alerts", alerts.Handler(alertManager, logger))

	if cfg.Metrics {
		api.logger.Info("enable exposing prometheus metrics")
		m.Handle("/metrics", promhttp.Handler())
		metrics.Register(logger)
	}

	api.server.Handler = m

	return api
}

func (api *API) Run(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	ln, err := net.Listen("tcp4", api.address)
	if err != nil {
		api.logger.Error("error listen address for api server", zap.String("address", api.address), zap.Error(err))
		ctxCancel()
		return
	}

	go func() {
		api.logger.Info("serve api server", zap.String("address", api.address))
		if err := api.server.Serve(ln); err != nil && err.Error() != "http: Server closed" {
			api.logger.Error("error serve api server", zap.Error(err))
			ctxCancel()
		}
	}()

	<-ctx.Done()

	api.logger.Info("shutdown api server")

	if err := api.server.Shutdown(ctx); err != nil {
		api.logger.Error("error shutdown api server", zap.Error(err))
	}
}
