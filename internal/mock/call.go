package mock

import (
	lua "github.com/yuin/gopher-lua"
)

func (m *ModuleMock) call(method string) lua.LGFunction {
	return func(L *lua.LState) int {

		var args []lua.LValue

		for i := 0; i < L.GetTop(); i++ {
			args = append(args, L.Get(i+1))
		}

		m.registry.AddQuery(method, args) // todo: check error

		resp, err := m.registry.Response(AnyValue, method, args)
		if err != nil {
			m.errors = append(m.errors, "error get response: "+err.Error())
			return 0
		}

		for _, a := range resp {
			L.Push(a)
		}

		return len(resp)
	}
}
