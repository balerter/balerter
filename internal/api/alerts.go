package api

import (
	"net/http"
)

func (api *API) handlerAlerts(rw http.ResponseWriter, req *http.Request) {
	//info := api.alertManager.GetAlerts()
	//
	//data, err := json.Marshal(info)
	//if err != nil {
	//	api.logger.Error("error marshaling response", zap.Error(err))
	//	rw.WriteHeader(500)
	//	return
	//}
	//
	//rw.Header().Add("Content-Type", "application/json")
	//
	//if _, err = fmt.Fprintf(rw, "%s", data); err != nil {
	//	api.logger.Error("error write response", zap.Error(err))
	//	rw.WriteHeader(500)
	//	return
	//}
}
