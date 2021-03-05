package alerts

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	queryArgLevels = "levels"
)

// GET /api/v1/alerts
//
// Endpoint receive arguments:
// level=error,success - filter by alert level. comma-separated
//
// Examples:
// GET /api/v1/alerts?level=error
// GET /api/v1/alerts?level=error,warn
func (a *Alerts) handlerIndex(rw http.ResponseWriter, req *http.Request) {
	var levels []alert.Level
	if s := req.URL.Query().Get(queryArgLevels); s != "" {
		for _, ls := range strings.Split(s, ",") {
			l, err := alert.LevelFromString(ls)
			if err != nil {
				http.Error(rw, fmt.Sprintf("error parse level %s, %v", ls, err), http.StatusBadRequest)
				return
			}
			levels = append(levels, l)
		}
	}

	data, err := a.alertManager.Index(levels)
	if err != nil {
		a.logger.Error("error get alerts index", zap.Error(err))
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}

	rw.Write(data.Marshal())
}
