package mock

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *ModuleMock) call(method string) lua.LGFunction {
	return func(L *lua.LState) int {

		var args []lua.LValue

		for i := 0; i < L.GetTop(); i++ {
			args = append(args, L.Get(i+1))
		}

		responseArgs, ok := m.responses[m.buildHash(method, args)]
		if !ok {
			m.logger.Error("unexpected call", zap.String("method name", method), zap.Any("args", args))
			m.errors = append(m.errors, fmt.Sprintf("unexpected method call %s with %v", method, args))
			return 0
		}

		for _, a := range responseArgs {
			L.Push(a)
		}

		m.queryLog[m.buildHash(method, args)]++

		return len(responseArgs)
	}
}
