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

		err := m.registry.AddQuery(method, args)
		if err != nil {
			err := "error add query: " + err.Error()
			m.logger.Error(err)
			m.errors = append(m.errors, err)
		}

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
