package registry

import (
	"github.com/balerter/balerter/internal/luaformatter"
	lua "github.com/yuin/gopher-lua"
)

func (r *Registry) Register(anyValue, method string, callArgs, retArgs []lua.LValue) error {
	e, ok := r.responseEntries[method]
	if !ok {
		e = newResponseEntry()
		r.responseEntries[method] = e
	}

	for _, a := range callArgs {
		s, err := luaformatter.ValueToString(a)
		if err != nil {
			return err
		}
		key := a.Type().String() + "@" + s
		if a.String() == anyValue {
			key = anyValue
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
