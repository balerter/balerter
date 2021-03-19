package prometheus

import (
	"fmt"
	"github.com/yuin/gluamapper"
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
	if options.Type() == lua.LTTable {
		err := gluamapper.Map(options.(*lua.LTable), &queryOptions)
		if err != nil {
			return nil, fmt.Errorf("error decode query options, %w", err)
		}
	}
	if err := queryOptions.validate(); err != nil {
		return nil, err
	}
	return queryOptions, nil
}
