package mock

import (
	"fmt"
	"github.com/balerter/balerter/internal/luaformatter"
	lua "github.com/yuin/gopher-lua"
)

func (m *ModuleMock) call(method string) lua.LGFunction {
	return func(luaState *lua.LState) int {

		var args []lua.LValue

		for i := 0; i < luaState.GetTop(); i++ {
			args = append(args, luaState.Get(i+1))
		}

		err := m.registry.AddCall(method, args)
		if err != nil {
			err := "error add query: " + err.Error()
			m.logger.Error(err)
			m.errors = append(m.errors, err)
		}

		resp, err := m.registry.Response(AnyValue, method, args)
		if err != nil {
			m.errors = append(m.errors, fmt.Sprintf("error get response for method '%s' with args '%s', %s",
				method, luaformatter.ValuesToStringNoErr(args), err.Error()))
			return 0
		}

		for _, a := range resp {
			luaState.Push(a)
		}

		return len(resp)
	}
}
