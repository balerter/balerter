package http

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type requestArgs struct {
	Method  string            `json:"method"`
	URI     string            `json:"uri"`
	Body    []byte            `json:"body"`
	Headers map[string]string `json:"headers"`
}

func newRequestArgs() *requestArgs {
	return &requestArgs{
		Headers: make(map[string]string),
	}
}

// parseRequestArgs parses arguments for methods http.get, http.post etc
// 1 - uri
// 2 - body
// 3 - headers
func (h *HTTP) parseRequestArgs(luaState *lua.LState) (*requestArgs, error) {
	args := newRequestArgs()

	argURI := luaState.Get(1)
	if argURI.Type() != lua.LTString {
		return nil, fmt.Errorf("uri argument must be a string")
	}
	args.URI = argURI.(lua.LString).String()

	argBody := luaState.Get(2)
	switch argBody.Type() {
	case lua.LTString:
		args.Body = []byte(argBody.(lua.LString).String())
	case lua.LTNil:
	default:
		return nil, fmt.Errorf("body argument must be a string or nil")
	}
	if argURI.Type() != lua.LTString {
		return nil, fmt.Errorf("first argument must be a string")
	}

	argHeaders := luaState.Get(3)
	switch argHeaders.Type() {
	case lua.LTTable:
		argHeaders.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			args.Headers[value.String()] = value2.String()
		})
	case lua.LTNil:
	default:
		return nil, fmt.Errorf("headers argument must be a table or nil")
	}

	return args, nil
}
