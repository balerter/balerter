package manager

import (
	"context"
	"errors"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/api/alerts"
	"github.com/balerter/balerter/internal/api/kv"
	"github.com/balerter/balerter/internal/api/runtime"
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

type httpServer interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}

type Runner interface {
	RunScript(name string, req *http.Request) error
}

type API struct {
	address string
	server  httpServer
	logger  *zap.Logger
}

func New(
	address string,
	coreStorageAlert,
	coreStorageKV coreStorage.CoreStorage,
	chManager ChManager,
	runner Runner,
	logger *zap.Logger,
) *API {
	alertsRouter := alerts.New(coreStorageAlert.Alert(), chManager, logger)
	kvRouter := kv.New(coreStorageKV.KV(), logger)
	runtimeRouter := runtime.New(runner, logger)

	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/alerts", alertsRouter.Handler)
		r.Route("/kv", kvRouter.Handler)
		r.Route("/runtime", runtimeRouter.Handler)
	})

	api := &API{
		address: address,
		server: &http.Server{
			Handler: router,
		},
		logger: logger,
	}

	return api
}

func (api *API) Run(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	go func() {
		api.logger.Info("serve api server", zap.String("address", api.address))
		e := api.server.Serve(ln)

		ctxCancel()

		if e != nil && !errors.Is(e, http.ErrServerClosed) {
			api.logger.Error("error serve api server", zap.Error(e))
		}
	}()

	<-ctx.Done()

	api.logger.Info("shutdown api server")

	err := api.server.Shutdown(context.Background())
	if err != nil {
		api.logger.Error("error shutdown api server", zap.Error(err))
	}
}
