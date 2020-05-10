package mock

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (m *ModuleMock) on(luaState *lua.LState) int {
	m.logger.Debug("[MOCK] on")

	if luaState.GetTop() == 0 {
		err := "mock.on should have first argument"
		m.logger.Error(err)
		m.errors = append(m.errors, err)
		return 0
	}

	methodNameL := luaState.Get(1)
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
	for i := 1; i < luaState.GetTop(); i++ {
		args = append(args, luaState.Get(i+1))
	}

	T := luaState.SetFuncs(luaState.NewTable(), map[string]lua.LGFunction{"response": m.saveResponse(methodName, args)})

	luaState.Push(T)

	return 1
}

func (m *ModuleMock) saveResponse(methodName string, callArgs []lua.LValue) lua.LGFunction {
	return func(luaState *lua.LState) int {
		retArgs := make([]lua.LValue, luaState.GetTop())
		for i := 0; i < luaState.GetTop(); i++ {
			retArgs[i] = luaState.Get(i + 1) // lua indexing starts with 1
		}

		err := m.registry.Register(AnyValue, methodName, callArgs, retArgs)
		if err != nil {
			m.errors = append(m.errors, "error register response: "+err.Error())
		}

		return 0
	}
}
