package log

import (
	"fmt"
	"net/http"
)

func (l *Log) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	switch method {
	case "error":
		l.logger.Error(string(body))
	case "warn":
		l.logger.Warn(string(body))
	case "info":
		l.logger.Info(string(body))
	case "debug":
		l.logger.Debug(string(body))
	default:
		return nil, http.StatusBadRequest, fmt.Errorf("unknod method %q", method)
	}

	return nil, http.StatusOK, nil
}
