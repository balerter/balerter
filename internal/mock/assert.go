package mock

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (m *ModuleMock) assert(called bool) lua.LGFunction {
	return func(L *lua.LState) int {
		if L.GetTop() == 0 {
			err := "assert should have method name as first argument"
			m.logger.Error(err)
			m.errors = append(m.errors, err)
			return 0
		}

		nameL := L.Get(1)
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
		for i := 1; i < L.GetTop(); i++ {
			args = append(args, L.Get(i+1))
		}

		hash := m.buildHash(name, args)

		if called {
			a, ok := m.assertsCalled[hash]
			if !ok {
				a = &assert{
					method: name,
					args:   args,
				}
				m.assertsCalled[hash] = a
			}
			a.count++
		} else {
			a, ok := m.assertsNotCalled[hash]
			if !ok {
				a = &assert{
					method: name,
					args:   args,
				}
				m.assertsNotCalled[hash] = a
			}
			a.count++
		}

		return 0
	}
}
