package test

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (m *Mock) onQuery(L *lua.LState) int {
	if L.GetTop() != 1 {
		m.errors = append(m.errors, fmt.Sprintf("wrong arguments count for 'onQuery' function: %d, expect 1", L.GetTop()))
		return 0
	}

	q := L.Get(1).String()

	T := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{"response": m.saveResponse(q)})

	L.Push(T)

	return 1
}

func (m *Mock) saveResponse(query string) lua.LGFunction {
	return func(L *lua.LState) int {
		args := make([]lua.LValue, L.GetTop())
		for i := 0; i < L.GetTop(); i++ {
			args[i] = L.Get(i + 1) // lua indexing starts with 1
		}

		m.responses[query] = args

		return 0
	}
}
