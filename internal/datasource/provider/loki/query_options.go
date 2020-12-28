package loki

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type queryOptions struct {
	Time      string
	Limit     int
	Direction string
}

func (q *queryOptions) validate() error {
	if err := directionValidate(q.Direction); err != nil {
		return err
	}
	return nil
}

func (m *Loki) parseQueryOptions(luaState *lua.LState) (*queryOptions, error) {
	queryOptions := &queryOptions{}

	options := luaState.Get(2)
	if options.Type() == lua.LTNil {
		return queryOptions, nil
	}

	if options.Type() != lua.LTTable {
		return nil, fmt.Errorf("options must be a table")
	}

	err := gluamapper.Map(options.(*lua.LTable), &queryOptions)
	if err != nil {
		return nil, fmt.Errorf("error parse, %w", err)
	}

	if err := queryOptions.validate(); err != nil {
		return nil, err
	}

	return queryOptions, nil
}
