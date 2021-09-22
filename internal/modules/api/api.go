package api

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
)

// ModuleName returns the module name
func ModuleName() string {
	return "api"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"is_api",
		"query",
		"url",
		"body",
		"host",
		"method",
	}
}

// API represents API core module
type API struct {
	isAPI  bool
	query  map[string][]string
	url    string
	body   string
	host   string
	method string
}

// New creates new API core module
func New() *API {
	a := &API{
		query: map[string][]string{},
	}

	return a
}

// Name returns the module name
func (a *API) Name() string {
	return ModuleName()
}

// FillData inits module data
func (a *API) FillData(req *http.Request) error {
	if req == nil {
		return nil
	}
	a.isAPI = true
	a.url = req.URL.String()
	a.method = req.Method
	a.host = req.Host

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	a.body = string(body)
	for key, values := range req.URL.Query() {
		if len(values) == 0 {
			return fmt.Errorf("no query values")
		}
		a.query[key] = values
	}

	return nil
}

// GetLoader returns the lua loader
func (a *API) GetLoader(_ modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"is_api": a.funcIsAPI,
				"query":  a.funcQuery,
				"url":    a.funcString(a.url),
				"body":   a.funcString(a.body),
				"host":   a.funcString(a.host),
				"method": a.funcString(a.method),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (a *API) Stop() error {
	return nil
}

func (a *API) funcIsAPI(luaState *lua.LState) int {
	luaState.Push(lua.LBool(a.isAPI))
	return 1
}

func (a *API) funcQuery(luaState *lua.LState) int {
	t := &lua.LTable{}
	for key, values := range a.query {
		tt := &lua.LTable{}
		for _, v := range values {
			tt.Append(lua.LString(v))
		}
		t.RawSetString(key, tt)
	}
	luaState.Push(t)
	return 1
}

func (a *API) funcString(v string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		luaState.Push(lua.LString(v))
		return 1
	}
}
