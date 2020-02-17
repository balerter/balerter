package api

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func (api *API) handlerConfig(rw http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(api.config)
	if err != nil {
		api.logger.Error("error marshaling response", zap.Error(err))
		rw.WriteHeader(500)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	if _, err = fmt.Fprintf(rw, "%s", data); err != nil {
		api.logger.Error("error write response", zap.Error(err))
		rw.WriteHeader(500)
		return
	}
}
