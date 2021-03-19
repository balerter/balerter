package runtime

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (rt *Runtime) handlerRun(rw http.ResponseWriter, req *http.Request) {
	scriptName := chi.URLParam(req, "name")
	if scriptName == "" {
		http.Error(rw, "empty name", http.StatusBadRequest)
		return
	}

	err := rt.runner.RunScript(scriptName, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
