package runtime

import (
	"fmt"
	"net/http"
)

func (m *Runtime) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method != "get" {
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
	}

	resp := struct {
		LogLevel     string `json:"log_level"`
		IsDebug      bool   `json:"is_debug"`
		IsOnce       bool   `json:"is_once"`
		WithScript   string `json:"with_script"`
		ConfigSource string `json:"config_source"`
	}{
		LogLevel:     m.flg.LogLevel,
		IsDebug:      m.flg.Debug,
		IsOnce:       m.flg.Once,
		WithScript:   m.flg.Script,
		ConfigSource: m.flg.ConfigFilePath,
	}

	return resp, http.StatusOK, nil
}
