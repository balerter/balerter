package prometheus

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type queryRangeOptions struct {
	Start string
	End   string
	Step  string
}

func (o *queryRangeOptions) validate() error {
	return nil
}

func (m *Prometheus) parseRangeOptions(luaState *lua.LState) (*queryRangeOptions, error) {
	options := luaState.Get(2)
	queryOptions := &queryRangeOptions{}
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
