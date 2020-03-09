package http

import lua "github.com/yuin/gopher-lua"

type response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

func newResponse() *response {
	return &response{
		Headers: make(map[string]string),
	}
}

func (r *response) toLuaTable() *lua.LTable {
	t := &lua.LTable{}
	t.RawSet(lua.LString("status_code"), lua.LNumber(r.StatusCode))
	t.RawSet(lua.LString("body"), lua.LString(r.Body))

	h := &lua.LTable{}
	for name, value := range r.Headers {
		h.RawSet(lua.LString(name), lua.LString(value))
	}

	t.RawSet(lua.LString("headers"), h)

	return t
}
