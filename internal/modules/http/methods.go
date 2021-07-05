package http

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (h *HTTP) send(method string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		args, err := h.parseRequestArgs(luaState)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error parse args, " + err.Error()))
			return 2
		}

		args.Method = method

		h.logger.Debug("http request", zap.Any("args", args))

		response, err := h.sendRequest(args)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error send request, " + err.Error()))
			return 2
		}

		luaState.Push(response.toLuaTable())

		return 1
	}
}

func (h *HTTP) request(luaState *lua.LState) int {
	opts := luaState.Get(1)
	if opts.Type() != lua.LTTable {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("argument must be a table"))
		return 2
	}

	args := newRequestArgs()

	err := args.parseFromTable(opts.(*lua.LTable))
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse arguments, " + err.Error()))
		return 2
	}

	h.logger.Debug("http request", zap.Any("args", args))

	response, err := h.sendRequest(args)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send request, " + err.Error()))
		return 2
	}

	luaState.Push(response.toLuaTable())

	return 1
}
