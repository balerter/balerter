package kv

import (
	"encoding/json"
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

	buf, err := json.Marshal(data)
	if err != nil {
		kv.logger.Error("error marshal kv data", zap.Error(err))
		http.Error(rw, "error marshal data", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(buf)
	if err != nil {
		kv.logger.Error("error write response", zap.Error(err))
		http.Error(rw, "error write response", http.StatusInternalServerError)
		return
	}
}
