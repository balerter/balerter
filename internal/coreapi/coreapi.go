package coreapi

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/modules"

	"go.uber.org/zap"
)

type coreapiHandlerFunc func(method string, parts []string, params map[string]string, body []byte) (any, int, error)

type CoreAPI struct {
	authToken string
	modules   map[string]coreapiHandlerFunc
	logger    *zap.Logger
}

func New(dsManager *manager.Manager, coreModules []modules.Module, authToken string, logger *zap.Logger) *CoreAPI {
	r := &CoreAPI{
		authToken: authToken,
		modules:   map[string]coreapiHandlerFunc{},
		logger:    logger,
	}
	for _, m := range coreModules {
		r.modules[m.Name()] = m.CoreApiHandler
	}
	r.modules["datasource"] = dsManager.CoreApiHandler

	return r
}

func (r *CoreAPI) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	srv := http.Server{
		Handler: http.HandlerFunc(r.handler),
	}

	go func() {
		<-ctx.Done()
		r.logger.Info("stop coreapi server")
		err := srv.Shutdown(context.Background())
		if err != nil {
			r.logger.Error("error shutdown coreapi server", zap.Error(err))
		}
	}()

	r.logger.Info("serve coreapi server", zap.String("address", ln.Addr().String()))
	errServe := srv.Serve(ln)
	if errServe != nil {
		if !errors.Is(errServe, http.ErrServerClosed) {
			r.logger.Error("error serve coreapi server", zap.Error(errServe))
		}
		cancel()
	}
}

type response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

func (r *CoreAPI) errorResponse(rw http.ResponseWriter, code int, err error) {
	resp := response{
		Status: "error",
		Error:  err.Error(),
	}

	respJson, errMarshal := json.Marshal(resp)
	if errMarshal != nil {
		r.logger.Error("error marshal response", zap.Error(errMarshal))
		http.Error(rw, "error marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(code)
	rw.Write(respJson)
}
