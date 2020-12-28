package loki

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type rangeOptions struct {
	Limit     int
	Start     string
	End       string
	Step      string
	Direction string
}

func (q *rangeOptions) validate() error {
	if err := directionValidate(q.Direction); err != nil {
		return err
	}
	return nil
}

func (m *Loki) parseRangeOptions(luaState *lua.LState) (*rangeOptions, error) {
	queryOptions := &rangeOptions{}

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
