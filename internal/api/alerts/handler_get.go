package alerts

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

func (a *Alerts) handlerGet(rw http.ResponseWriter, req *http.Request) {
	alertName := chi.URLParam(req, "name")
	if alertName == "" {
		http.Error(rw, "empty name", http.StatusBadRequest)
		return
	}

	alert, err := a.alertManager.Get(alertName)
	if err != nil {
		a.logger.Error("error get alert", zap.Error(err))
		http.Error(rw, "error get alert", http.StatusInternalServerError)
		return
	}

	if alert == nil {
		http.Error(rw, "alert not found", http.StatusNotFound)
		return
	}

	rw.Write(alert.Marshal())
}
