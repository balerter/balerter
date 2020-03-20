package alerts

import (
	coreStorage "github.com/balerter/balerter/internal/core_storage"
	"go.uber.org/zap"
	"net/http"
)

// HandlerIndex handle API request GET /api/v1/alerts
//
// Endpoint receive arguments:
// name=<NAME1>,<NAME2> - filter by name
// level=error,success - filter by alert level
//
// Examples:
// GET /api/v1/alerts?level=error
// GET /api/v1/alerts?level=error,warn&name=foo
// GET /api/v1/alerts?level=error,warn&name=foo,bar
func HandlerIndex(coreStorageAlert coreStorage.CoreStorage, logger *zap.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		data, err := coreStorageAlert.Alert().All()
		if err != nil {
			logger.Error("error get alerts", zap.Error(err))
			rw.Header().Add("X-Error", err.Error())
			rw.WriteHeader(500)
			return
		}

		data, err = filter(req, data)
		if err != nil {
			logger.Error("error filter alerts", zap.Error(err))
			rw.Header().Add("X-Error", err.Error())
			rw.WriteHeader(400)
			return
		}

		err = newResource(data).render(rw)
		if err != nil {
			logger.Error("error write response", zap.Error(err))
			rw.Header().Add("X-Error", "error write response")
			rw.WriteHeader(500)
			return
		}
	}
}
