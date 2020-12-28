package kv

import (
	"go.uber.org/zap"
	"net/http"
)

// GET /api/v1/kv
func (kv *KV) handlerIndex(rw http.ResponseWriter, _ *http.Request) {
	var err error

	data, err := kv.storage.All()
	if err != nil {
		kv.logger.Error("error get kv data", zap.Error(err))
		rw.Header().Add("X-Error", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = newResource(data).render(rw)
	if err != nil {
		kv.logger.Error("error write response", zap.Error(err))
		rw.Header().Add("X-Error", "error write response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
