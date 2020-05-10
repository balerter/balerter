package manager

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/api/alerts"
	"github.com/balerter/balerter/internal/api/kv"
	"github.com/balerter/balerter/internal/config"
	coreStorage "github.com/balerter/balerter/internal/core_storage"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"
)

const (
	pprofPrefix = "/debug/pprof"
)

type API struct {
	address string
	server  *http.Server
	logger  *zap.Logger
}

func New(cfg config.API, coreStorageAlert, coreStorageKV coreStorage.CoreStorage, logger *zap.Logger) *API {
	api := &API{
		address: cfg.Address,
		server:  &http.Server{},
		logger:  logger,
	}

	m := http.NewServeMux()

	m.HandleFunc("/liveness", api.handlerLiveness)
	m.HandleFunc("/api/v1/alerts", alerts.HandlerIndex(coreStorageAlert, logger))
	m.HandleFunc("/api/v1/kv", kv.HandlerIndex(coreStorageKV, logger))

	m.HandleFunc(pprofPrefix+"/profile", pprof.Profile)
	m.HandleFunc(pprofPrefix+"/trace", pprof.Trace)
	m.HandleFunc(pprofPrefix+"/heap", pprof.Handler("heap").ServeHTTP)
	m.HandleFunc(pprofPrefix+"/goroutine", pprof.Handler("goroutine").ServeHTTP)
	m.HandleFunc(pprofPrefix+"/allocs", pprof.Handler("allocs").ServeHTTP)

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

func (api *API) handlerLiveness(rw http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprint(rw, "ok"); err != nil {
		api.logger.Error("error write response", zap.Error(err))
		rw.WriteHeader(500)
	}
}
