package coreapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/modules"

	"go.uber.org/zap"
)

type coreapiHandlerFunc func(req []string, body []byte) (any, int, error)

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

func (r *CoreAPI) handler(rw http.ResponseWriter, req *http.Request) {
	if r.authToken != "" {
		if req.Header.Get("Authorization") != r.authToken {
			r.errorResponse(rw, http.StatusForbidden, fmt.Errorf("forbidden"))
			return
		}
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) == 0 {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("empty path"))
		return
	}

	moduleName := parts[0]
	parts = parts[1:]

	h, ok := r.modules[moduleName]
	if !ok {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("module %q not found", moduleName))
		return
	}

	reqBody, errReqBody := io.ReadAll(req.Body)
	if errReqBody != nil {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("error read request body, %s", errReqBody))
		return
	}

	resp, errCode, err := h(parts, reqBody)
	if err != nil {
		r.logger.Error("error call coreapi handler", zap.String("module", moduleName), zap.Int("errCode", errCode), zap.Error(err))
		r.errorResponse(rw, errCode, err)
		return
	}

	rsp := response{
		Status: "ok",
		Result: resp,
	}

	respJson, errMarshal := json.Marshal(rsp)
	if errMarshal != nil {
		r.logger.Error("error marshal response", zap.Error(errMarshal))
		r.errorResponse(rw, http.StatusInternalServerError, fmt.Errorf("error marshal response"))
		return
	}

	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(respJson)
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
