package api

import (
	"context"
	"fmt"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

type alertManagerAPIer interface {
	GetAlerts() []*alertManager.APIAlertInfo
}

type API struct {
	address      string
	server       *http.Server
	alertManager alertManagerAPIer
	logger       *zap.Logger
	config       config.Config

	pathPrefix string
}

func New(cfg config.API, fullConfig config.Config, alertManager alertManagerAPIer, logger *zap.Logger) *API {
	api := &API{
		address:      cfg.Address,
		server:       &http.Server{},
		alertManager: alertManager,
		logger:       logger,
		pathPrefix:   fmt.Sprintf("/api/%s", "v1"),
		config:       fullConfig,
	}

	mux := http.NewServeMux()

	mux.HandleFunc(api.pathPrefix+"/alerts", api.handlerAlerts)
	mux.HandleFunc(api.pathPrefix+"/config", api.handlerConfig)

	api.server.Handler = mux

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
