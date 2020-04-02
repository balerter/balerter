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

func ModuleName() string {
	return "http"
}

func Methods() []string {
	return []string{
		"request",
		"post",
		"get",
		"put",
		"delete",
	}
}

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
	return ModuleName()
}

func (h *HTTP) GetLoader(_ *script.Script) lua.LGFunction {
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

			mod.RawSetString("methodGet", lua.LString(http.MethodGet))
			mod.RawSetString("methodHead", lua.LString(http.MethodHead))
			mod.RawSetString("methodPost", lua.LString(http.MethodPost))
			mod.RawSetString("methodPut", lua.LString(http.MethodPut))
			mod.RawSetString("methodPatch", lua.LString(http.MethodPatch))
			mod.RawSetString("methodDelete", lua.LString(http.MethodDelete))
			mod.RawSetString("methodConnect", lua.LString(http.MethodConnect))
			mod.RawSetString("methodOptions", lua.LString(http.MethodOptions))
			mod.RawSetString("methodTrace", lua.LString(http.MethodTrace))

			L.Push(mod)
			return 1
		}
	}()
}

func (h *HTTP) Stop() error {
	h.client.CloseIdleConnections()

	return nil
}
