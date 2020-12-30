package alerts

import (
	"encoding/json"
	"github.com/balerter/balerter/internal/alert"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type alertUpdatePayload struct {
	Level string `json:"level"`
	Text  string `json:"text"`
	// Deprecated
	Fields   []string `json:"fields,omitempty"`
	Channels []string `json:"channels,omitempty"`
	Quiet    bool     `json:"quiet,omitempty"`
	Repeat   int      `json:"repeat,omitempty"`
	Image    string   `json:"image,omitempty"`
}

func (a *Alerts) handlerUpdate(rw http.ResponseWriter, req *http.Request) {
	alertName := chi.URLParam(req, "name")
	if alertName == "" {
		http.Error(rw, "empty name", 400)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		a.logger.Error("error read body", zap.Error(err))
		http.Error(rw, "error read body", 400)
		return
	}
	defer req.Body.Close()

	p := &alertUpdatePayload{}

	err = json.Unmarshal(buf, p)
	if err != nil {
		a.logger.Error("error unmarshal body", zap.Error(err))
		http.Error(rw, "error unmarshal body", 400)
		return
	}

	l, err := alert.LevelFromString(p.Level)
	if err != nil {
		a.logger.Error("error parse level", zap.Error(err))
		http.Error(rw, "error parse level", 400)
		return
	}

	_ = l
	// TODO: implement it
	//opts := &alert.Options{
	//	Fields:   p.Fields,
	//	Channels: p.Channels,
	//	Quiet:    p.Quiet,
	//	Repeat:   p.Repeat,
	//	Image:    p.Image,
	//}
	//err = a.alertManager.Update(alertName, l, p.Text, opts)
	//if err != nil {
	//	a.logger.Error("error update alert", zap.Error(err))
	//	http.Error(rw, "error update alert", 400)
	//	return
	//}
}
