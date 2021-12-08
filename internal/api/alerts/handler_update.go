package alerts

import (
	"encoding/json"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type alertUpdatePayload struct {
	Level    string   `json:"level"`
	Text     string   `json:"text"`
	Channels []string `json:"channels,omitempty"`
	Quiet    bool     `json:"quiet,omitempty"`
	Repeat   int      `json:"repeat,omitempty"`
	Image    string   `json:"image,omitempty"`
}

func (a *Alerts) handlerUpdate(rw http.ResponseWriter, req *http.Request) {
	alertName := chi.URLParam(req, "name")
	if alertName == "" {
		http.Error(rw, "empty name", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		a.logger.Error("error read body", zap.Error(err))
		http.Error(rw, "error read body", http.StatusInternalServerError)
		return
	}

	payload := &alertUpdatePayload{}

	err = json.Unmarshal(buf, payload)
	if err != nil {
		a.logger.Error("error unmarshal body", zap.Error(err))
		http.Error(rw, fmt.Sprintf("error unmarshal body, %v", err), http.StatusBadRequest)
		return
	}

	l, err := alert.LevelFromString(payload.Level)
	if err != nil {
		http.Error(rw, fmt.Sprintf("error parse level %s, %v", payload.Level, err), http.StatusBadRequest)
		return
	}

	updatedAlert, levelWasUpdated, err := a.alertManager.Update(alertName, l)
	if err != nil {
		a.logger.Error("error update alert", zap.Error(err))
		http.Error(rw, "error update alert", http.StatusInternalServerError)
		return
	}

	if levelWasUpdated || (payload.Repeat > 0 && updatedAlert.Count%payload.Repeat == 0) {
		a.chManager.Send(updatedAlert, payload.Text, &alert.Options{
			Channels: payload.Channels,
			Quiet:    payload.Quiet,
			Repeat:   payload.Repeat,
			Image:    payload.Image,
		})
	}

	rw.Write(updatedAlert.Marshal())
}
