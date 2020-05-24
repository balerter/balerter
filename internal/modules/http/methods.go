package http

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (h *HTTP) send(method string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		args, err := h.parseRequestArgs(luaState)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error parse args, " + err.Error()))
			return 2 //nolint:mnd
		}

		args.Method = method

		h.logger.Debug("http request", zap.Any("args", args))

		response, err := h.sendRequest(args)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error send request, " + err.Error()))
			return 2 //nolint:mnd
		}

		luaState.Push(response.toLuaTable())

		return 1 //nolint:mnd
	}
}

func (h *HTTP) request(luaState *lua.LState) int {
	opts := luaState.Get(1) //nolint:mnd
	if opts.Type() != lua.LTTable {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("argument must be a table"))
		return 2 //nolint:mnd
	}

	args := newRequestArgs()

	err := gluamapper.Map(opts.(*lua.LTable), args)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse arguments, " + err.Error()))
		return 2 //nolint:mnd
	}

	h.logger.Debug("http request", zap.Any("args", args))

	response, err := h.sendRequest(args)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send request, " + err.Error()))
		return 2 //nolint:mnd
	}

	luaState.Push(response.toLuaTable())

	return 1 //nolint:mnd
}
