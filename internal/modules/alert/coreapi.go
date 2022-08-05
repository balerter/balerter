package alert

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/balerter/balerter/internal/alert"

	"go.uber.org/zap"
)

type coreapiRequest struct {
	Message  string           `json:"message"`
	Options  alert.Options    `json:"options"`
	Escalate map[int][]string `json:"escalate"`
}

type coreapiResponse struct {
	Alert           *alert.Alert `json:"alert"`
	LevelWasUpdated bool         `json:"level_was_updated"`
}

// CoreApiHandler expect req:
// /alert/<method>/<name>

func (a *Alert) CoreApiHandler(req []string, body []byte) (any, int, error) {
	if len(req) != 2 {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request, expected: /alert/{method}/{name}")
	}

	name := req[1]

	var level alert.Level

	switch req[0] {
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
		return nil, http.StatusBadRequest, fmt.Errorf("invalid method: %s", req[0])
	}

	r := coreapiRequest{}

	errUnmarshalBody := json.Unmarshal(body, &r)
	if errUnmarshalBody != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error unmarshal body: %s", errUnmarshalBody)
	}

	if r.Message == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("message is required")
	}

	updatedAlert, levelWasUpdated, err := a.call(name, r.Message, nil, r.Escalate, level, &r.Options)
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
