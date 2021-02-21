package manager

import (
	"context"
	"errors"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/api/alerts"
	"github.com/balerter/balerter/internal/api/kv"
	apiConfig "github.com/balerter/balerter/internal/config/global/api"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

type ChManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

type API struct {
	address string
	server  *http.Server
	logger  *zap.Logger
}

func New(cfg apiConfig.API, coreStorageAlert, coreStorageKV coreStorage.CoreStorage, chManager ChManager, logger *zap.Logger) *API {
	api := &API{
		address: cfg.Address,
		server:  &http.Server{},
		logger:  logger,
	}

	alertsRouter := alerts.New(coreStorageAlert.Alert(), chManager, logger)
	kvRouter := kv.New(coreStorageKV.KV(), logger)

	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/alerts", alertsRouter.Handler)
		r.Route("/kv", kvRouter.Handler)
	})

	api.server.Handler = router

	return api
}

func (api *API) Run(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	go func() {
		api.logger.Info("serve api server", zap.String("address", api.address))
		e := api.server.Serve(ln)

		if e != nil && !errors.Is(e, http.ErrServerClosed) {
			api.logger.Error("error serve api server", zap.Error(e))
			ctxCancel()
		}
	}()

	<-ctx.Done()

	api.logger.Info("shutdown api server")

	err := api.server.Shutdown(context.Background())
	if err != nil {
		api.logger.Error("error shutdown api server", zap.Error(err))
	}
}
