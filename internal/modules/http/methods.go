package http

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (h *HTTP) send(method string) lua.LGFunction {
	return func(L *lua.LState) int {
		args, err := h.parseRequestArgs(L)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("error parse args, " + err.Error()))
			return 2
		}

		args.Method = method

		h.logger.Debug("http request", zap.Any("args", args))

		response, err := h.sendRequest(args)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("error send request, " + err.Error()))
			return 2
		}

		L.Push(response.toLuaTable())

		return 1
	}
}

func (h *HTTP) request(L *lua.LState) int {
	opts := L.Get(1)
	if opts.Type() != lua.LTTable {
		L.Push(lua.LNil)
		L.Push(lua.LString("argument must be a table"))
		return 2
	}

	args := newRequestArgs()

	err := gluamapper.Map(opts.(*lua.LTable), args)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error parse arguments, " + err.Error()))
		return 2
	}

	h.logger.Debug("http request", zap.Any("args", args))

	response, err := h.sendRequest(args)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error send request, " + err.Error()))
		return 2
	}

	L.Push(response.toLuaTable())

	return 1
}
