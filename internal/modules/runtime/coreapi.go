package runtime

import (
	"net/http"
)

func (m *Runtime) CoreApiHandler(_ []string, _ []byte) (any, int, error) {
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
