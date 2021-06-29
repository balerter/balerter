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

// ModuleName returns the module name
func ModuleName() string {
	return "http"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"request",
		"post",
		"get",
		"put",
		"delete",
	}
}

// HTTP represents the HTTP core module
type HTTP struct {
	logger *zap.Logger
	client *http.Client
}

// New creates HTTP core module
func New(logger *zap.Logger) *HTTP {
	h := &HTTP{
		logger: logger,
	}

	h.client = &http.Client{
		Timeout: defaultTimeout,
	}

	return h
}

// Name returns the module name
func (h *HTTP) Name() string {
	return ModuleName()
}

// GetLoader returns the lua loader
func (h *HTTP) GetLoader(_ *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"request": h.request,
				"post":    h.send(http.MethodPost),
				"get":     h.send(http.MethodGet),
				"put":     h.send(http.MethodPut),
				"delete":  h.send(http.MethodDelete),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			mod.RawSetString("methodGet", lua.LString(http.MethodGet))
			mod.RawSetString("methodHead", lua.LString(http.MethodHead))
			mod.RawSetString("methodPost", lua.LString(http.MethodPost))
			mod.RawSetString("methodPut", lua.LString(http.MethodPut))
			mod.RawSetString("methodPatch", lua.LString(http.MethodPatch))
			mod.RawSetString("methodDelete", lua.LString(http.MethodDelete))
			mod.RawSetString("methodConnect", lua.LString(http.MethodConnect))
			mod.RawSetString("methodOptions", lua.LString(http.MethodOptions))
			mod.RawSetString("methodTrace", lua.LString(http.MethodTrace))

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (h *HTTP) Stop() error {
	h.client.CloseIdleConnections()

	return nil
}
