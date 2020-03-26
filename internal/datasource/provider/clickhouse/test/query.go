package test

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (m *Mock) query(L *lua.LState) int {
	if L.GetTop() != 1 {
		m.errors = append(m.errors, fmt.Sprintf("wrong arguments count for 'query' function: %d, expect 1", L.GetTop()))
		return 0
	}

	q := L.Get(1).String()

	m.queryLog[q]++

	args, ok := m.responses[q]
	if !ok {
		m.errors = append(m.errors, fmt.Sprintf("a response is not defined for the query: %s", q))
		return 0
	}

	for _, a := range args {
		L.Push(a)
	}

	return len(args)
}
