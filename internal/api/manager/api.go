package manager

import (
	"context"
	"github.com/balerter/balerter/internal/api/alerts"
	"github.com/balerter/balerter/internal/api/kv"
	apiConfig "github.com/balerter/balerter/internal/config/global/api"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

type API struct {
	address string
	server  *http.Server
	logger  *zap.Logger
}

func New(cfg apiConfig.API, coreStorageAlert, coreStorageKV coreStorage.CoreStorage, logger *zap.Logger) *API {
	api := &API{
		address: cfg.Address,
		server:  &http.Server{},
		logger:  logger,
	}

	m := http.NewServeMux()

	m.HandleFunc("/api/v1/alerts", alerts.HandlerIndex(coreStorageAlert, logger))
	m.HandleFunc("/api/v1/kv", kv.HandlerIndex(coreStorageKV, logger))

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
		e := api.server.Serve(ln)
		if e != nil && e.Error() != "http: Server closed" {
			api.logger.Error("error serve api server", zap.Error(e))
			ctxCancel()
		}
	}()

	<-ctx.Done()

	api.logger.Info("shutdown api server")

	err = api.server.Shutdown(ctx)
	if err != nil {
		if err.Error() != "context canceled" {
			api.logger.Error("error shutdown api server", zap.Error(err))
		}
	}
}
