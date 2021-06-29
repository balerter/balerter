package service

import (
	"context"
	"errors"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"
)

// Service represents the Service module
type Service struct {
	server *http.Server
	logger *zap.Logger
}

var (
	livenessResponse = []byte("ok")
)

// New creates new Service
func New(logger *zap.Logger) *Service {
	s := &Service{
		server: &http.Server{},
		logger: logger,
	}

	metrics.Register(logger)

	router := chi.NewRouter()
	router.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/profile", pprof.Profile)
		r.Get("/trace", pprof.Trace)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
	})

	router.Get("/liveness", s.livenessHandler)
	router.Get("/metrics", promhttp.Handler().ServeHTTP)

	s.server.Handler = router

	return s
}

// Run the module
func (s *Service) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	go func() {
		s.logger.Info("serve service server", zap.String("address", ln.Addr().String()))
		err := s.server.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("error serve service server", zap.Error(err))
			cancel()
		}
	}()

	<-ctx.Done()

	s.logger.Info("shutdown service server")

	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.logger.Error("error shutdown service server", zap.Error(err))
	}
}

func (s *Service) livenessHandler(rw http.ResponseWriter, _ *http.Request) {
	_, err := rw.Write(livenessResponse)
	if err != nil {
		s.logger.Error("error write http response", zap.Error(err))
	}
}
