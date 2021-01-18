package alerts

import (
	"net/http"
)

// GET /api/v1/alerts
//
// Endpoint receive arguments:
// name=<NAME1>,<NAME2> - filter by name
// level=error,success - filter by alert level
//
// Examples:
// GET /api/v1/alerts?level=error
// GET /api/v1/alerts?level=error,warn&name=foo
// GET /api/v1/alerts?level=error,warn&name=foo,bar
func (a *Alerts) handlerIndex(rw http.ResponseWriter, req *http.Request) {
	// TODO: wip
	//var err error
	//
	//data, err := a.alertManager.All()
	//if err != nil {
	//	a.logger.Error("error get alerts", zap.Error(err))
	//	rw.Header().Add("X-Error", err.Error())
	//	rw.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//data, err = filter(req, data)
	//if err != nil {
	//	a.logger.Error("error filter alerts", zap.Error(err))
	//	rw.Header().Add("X-Error", err.Error())
	//	rw.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//err = newResource(data).render(rw)
	//if err != nil {
	//	a.logger.Error("error write response", zap.Error(err))
	//	rw.Header().Add("X-Error", "error write response")
	//	rw.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
}
