package registry

import (
	lua "github.com/yuin/gopher-lua"
)

func (r *Registry) AddCall(method string, args []lua.LValue) error {

	r.calls = append(r.calls, call{
		method: method,
		args:   args,
	})

	return nil
}
