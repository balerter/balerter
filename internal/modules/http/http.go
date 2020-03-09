package http

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Second * 30
)

type HTTP struct {
	logger *zap.Logger
	client *http.Client
}

func New(logger *zap.Logger) *HTTP {
	h := &HTTP{
		logger: logger,
	}

	h.client = &http.Client{
		Timeout: defaultTimeout,
	}

	return h
}

func (h *HTTP) Name() string {
	return "http"
}

func (h *HTTP) GetLoader(script *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(L *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"request": h.request,
				"post":    h.send(http.MethodPost),
				"get":     h.send(http.MethodGet),
				"put":     h.send(http.MethodPut),
				"delete":  h.send(http.MethodDelete),
			}

			mod := L.SetFuncs(L.NewTable(), exports)

			mod.RawSetString("methodGet", lua.LString("GET"))
			mod.RawSetString("methodHead", lua.LString("HEAD"))
			mod.RawSetString("methodPost", lua.LString("POST"))
			mod.RawSetString("methodPut", lua.LString("PUT"))
			mod.RawSetString("methodPatch", lua.LString("PATCH"))
			mod.RawSetString("methodDelete", lua.LString("DELETE"))
			mod.RawSetString("methodConnect", lua.LString("CONNECT"))
			mod.RawSetString("methodOptions", lua.LString("OPTIONS"))
			mod.RawSetString("methodTrace", lua.LString("TRACE"))

			L.Push(mod)
			return 1
		}
	}()
}

func (h *HTTP) Stop() error {
	h.client.CloseIdleConnections()

	return nil
}
