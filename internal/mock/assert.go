package mock

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (m *ModuleMock) assert(called bool) lua.LGFunction {
	return func(luaState *lua.LState) int {
		if luaState.GetTop() == 0 {
			err := "query should have a method name as first argument"
			m.logger.Error(err)
			m.errors = append(m.errors, err)
			return 0
		}

		nameL := luaState.Get(1)
		if nameL.Type() != lua.LTString {
			err := "method name be a string"
			m.logger.Error(err)
			m.errors = append(m.errors, err)
			return 0
		}

		name := strings.TrimSpace(nameL.String())
		if name == "" {
			err := "method name be a not empty"
			m.logger.Error(err)
			m.errors = append(m.errors, err)
			return 0
		}

		var args []lua.LValue
		for i := 1; i < luaState.GetTop(); i++ {
			args = append(args, luaState.Get(i+1))
		}

		err := m.registry.AddAssert(name, args, called)
		if err != nil {
			err := "error register assert: " + err.Error()
			m.logger.Error(err)
			m.errors = append(m.errors, err)
		}

		return 0
	}
}
