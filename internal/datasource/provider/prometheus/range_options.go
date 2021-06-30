package prometheus

import (
	"fmt"
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
	rangeOptions := &queryRangeOptions{}

	if options.Type() == lua.LTNil {
		return rangeOptions, nil
	}

	if options.Type() != lua.LTTable {
		return nil, fmt.Errorf("options must be a table")
	}

	optionsTable := options.(*lua.LTable)

	v := optionsTable.RawGetString("start")
	if v.Type() != lua.LTNil {
		if v.Type() != lua.LTString {
			return nil, fmt.Errorf("start must be a string")
		}
		rangeOptions.Start = string(v.(lua.LString))
	}

	v = optionsTable.RawGetString("end")
	if v.Type() != lua.LTNil {
		if v.Type() != lua.LTString {
			return nil, fmt.Errorf("end must be a string")
		}
		rangeOptions.End = string(v.(lua.LString))
	}

	v = optionsTable.RawGetString("step")
	if v.Type() != lua.LTNil {
		if v.Type() != lua.LTString {
			return nil, fmt.Errorf("step must be a string")
		}
		rangeOptions.Step = string(v.(lua.LString))
	}

	if err := rangeOptions.validate(); err != nil {
		return nil, err
	}
	return rangeOptions, nil
}
