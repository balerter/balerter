package alert

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/balerter/balerter/internal/alert"

	"go.uber.org/zap"
)

type coreapiResponse struct {
	Alert           *alert.Alert `json:"alert"`
	LevelWasUpdated bool         `json:"level_was_updated"`
}

func (a *Alert) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if len(parts) != 1 {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request, expected: /alert/{method}/{name}")
	}

	name := parts[0]

	var level alert.Level

	switch method {
	case "warn", "warning":
		level = alert.LevelWarn
	case "error", "fail":
		level = alert.LevelError
	case "success", "ok":
		level = alert.LevelSuccess
	case "get":
		al, errGetAlert := a.storage.Get(name)
		if errGetAlert != nil {
			return nil, http.StatusInternalServerError, errGetAlert
		}
		if al == nil {
			return nil, http.StatusOK, fmt.Errorf("alert not found")
		}
		return al, 0, nil
	default:
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method: %s", method)
	}

	if len(body) == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("empty body")
	}

	escalate := map[int][]string{}

	var opts = alert.Options{}

	if v, ok := params["channels"]; ok {
		opts.Channels = strings.Split(v, ",")
	}
	if v, ok := params["quiet"]; ok {
		if v == "true" {
			opts.Quiet = true
		}
	}
	if v, ok := params["repeat"]; ok {
		vv, err := strconv.Atoi(v)
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid repeat value, %v", v)
		}
		opts.Repeat = vv
	}
	if v, ok := params["image"]; ok {
		opts.Image = v
	}
	if v, ok := params["fields"]; ok {
		for _, s := range strings.Split(v, ",") {
			p := strings.Split(s, ":")
			if len(p) != 2 {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid fields value: %s", s)
			}
			opts.Fields[p[0]] = p[1]
		}
	}
	if v, ok := params["escalate"]; ok {
		for _, s := range strings.Split(v, ";") {
			p := strings.Split(s, ":")
			if len(p) != 2 {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid escalate value: %s", s)
			}
			n, err := strconv.Atoi(p[0])
			if err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid escalate value: %s", s)
			}
			escalate[n] = strings.Split(p[1], ",")
		}
	}

	updatedAlert, levelWasUpdated, err := a.call(name, string(body), nil, escalate, level, &opts)
	if err != nil {
		a.logger.Error("error alert.call", zap.Error(err))
		return nil, http.StatusInternalServerError, fmt.Errorf("internal error")
	}

	resp := coreapiResponse{
		Alert:           updatedAlert,
		LevelWasUpdated: levelWasUpdated,
	}

	return resp, 0, nil
}
