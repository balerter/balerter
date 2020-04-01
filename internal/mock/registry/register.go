package registry

import (
	"github.com/balerter/balerter/internal/lua_formatter"
	lua "github.com/yuin/gopher-lua"
)

func (r *Registry) Register(AnyValue, method string, callArgs, retArgs []lua.LValue) error {
	e, ok := r.responseEntries[method]
	if !ok {
		e = newResponseEntry()
		r.responseEntries[method] = e
	}

	for _, a := range callArgs {
		s, err := lua_formatter.ValueToString(a)
		if err != nil {
			return err
		}
		key := a.Type().String() + "@" + s
		if a.String() == AnyValue {
			key = AnyValue
		}

		nextEntry, ok := e.entries[key]
		if !ok {
			nextEntry = newResponseEntry()
			e.entries[key] = nextEntry
		}
		e = nextEntry
	}

	e.responses = append(e.responses, retArgs)

	return nil
}
