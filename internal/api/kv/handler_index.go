package kv

import (
	coreStorage "github.com/balerter/balerter/internal/core_storage"
	"go.uber.org/zap"
	"net/http"
)

// HandlerIndex handle API request GET /api/v1/kv
func HandlerIndex(coreStorage coreStorage.CoreStorage, logger *zap.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		data, err := coreStorage.KV().All()
		if err != nil {
			logger.Error("error get kv data", zap.Error(err))
			rw.Header().Add("X-Error", err.Error())
			rw.WriteHeader(500)
			return
		}

		err = newResource(data).render(rw)
		if err != nil {
			logger.Error("error write response", zap.Error(err))
			rw.Header().Add("X-Error", "error write response")
			rw.WriteHeader(500)
			return
		}
	}
}
