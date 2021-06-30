package prometheus

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type queryQueryOptions struct {
	Time string
}

func (q *queryQueryOptions) validate() error {
	return nil
}

func (m *Prometheus) parseQueryOptions(luaState *lua.LState) (*queryQueryOptions, error) {
	options := luaState.Get(2)
	queryOptions := &queryQueryOptions{}

	if options.Type() == lua.LTNil {
		return queryOptions, nil
	}
	if options.Type() != lua.LTTable {
		return nil, fmt.Errorf("options must be a table")
	}

	t := options.(*lua.LTable).RawGetString("time")
	if t.Type() == lua.LTNil {
		return queryOptions, nil
	}

	if t.Type() != lua.LTString {
		return nil, fmt.Errorf("time must be a string")
	}

	queryOptions.Time = string(t.(lua.LString))

	if err := queryOptions.validate(); err != nil {
		return nil, err
	}
	return queryOptions, nil
}
