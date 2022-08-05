package coreapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

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
	if len(parts) < 2 {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("invalid path"))
		return
	}
	moduleName := parts[0]
	method := parts[1]
	parts = parts[2:]

	params := map[string]string{}
	for k, v := range req.URL.Query() {
		params[k] = v[0]
	}

	reqBody, errReqBody := io.ReadAll(req.Body)
	if errReqBody != nil {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("error read request body, %s", errReqBody))
		return
	}

	h, ok := r.modules[moduleName]
	if !ok {
		r.errorResponse(rw, http.StatusBadRequest, fmt.Errorf("module %q not found", moduleName))
		return
	}

	resp, errCode, err := h(method, parts, params, reqBody)
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
