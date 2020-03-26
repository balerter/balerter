package test

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (m *Mock) assertCalledQuery(L *lua.LState) int {
	if L.GetTop() != 1 {
		m.errors = append(m.errors, fmt.Sprintf("wrong arguments count for 'assertCalledQuery' function: %d, expect 1", L.GetTop()))
		return 0
	}

	q := L.Get(1).String()

	m.asserts.callQuery[q]++

	return 0
}
