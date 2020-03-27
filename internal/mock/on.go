package mock

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

func (m *ModuleMock) on(L *lua.LState) int {
	m.logger.Debug("[MOCK] on")

	if L.GetTop() == 0 {
		err := "mock.on should have first argument"
		m.logger.Error(err)
		m.errors = append(m.errors, err)
		return 0
	}

	methodNameL := L.Get(1)
	if methodNameL.Type() != lua.LTString {
		err := "mock.on first argument should be a string"
		m.logger.Error(err)
		m.errors = append(m.errors, err)
		return 0
	}

	methodName := strings.TrimSpace(methodNameL.String())
	if methodName == "" {
		err := "mock.on first argument should be not empty"
		m.logger.Error(err)
		m.errors = append(m.errors, err)
		return 0
	}

	var args []lua.LValue
	for i := 1; i < L.GetTop(); i++ {
		args = append(args, L.Get(i+1))
	}

	T := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{"response": m.saveResponse(methodName, args)})

	L.Push(T)

	return 1
}

func (m *ModuleMock) saveResponse(methodName string, callArgs []lua.LValue) lua.LGFunction {
	return func(L *lua.LState) int {
		retArgs := make([]lua.LValue, L.GetTop())
		for i := 0; i < L.GetTop(); i++ {
			retArgs[i] = L.Get(i + 1) // lua indexing starts with 1
		}

		if _, ok := m.responses[m.buildHash(methodName, callArgs)]; ok {
			m.logger.Error("response value already defined for this request", zap.String("method name", methodName), zap.Any("values", callArgs))
			m.errors = append(m.errors, "response value already defined for this request, method "+methodName)
			return 0
		}

		m.responses[m.buildHash(methodName, callArgs)] = retArgs

		return 0
	}
}
