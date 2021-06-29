package registry

import (
	"github.com/balerter/balerter/internal/luaformatter"
	lua "github.com/yuin/gopher-lua"
)

// AddAssert add assert for test
func (r *Registry) AddAssert(method string, args []lua.LValue, called bool) error {
	e, ok := r.assertEntries[method]
	if !ok {
		e = newAssertEntry()
		r.assertEntries[method] = e
	}

	for _, a := range args {
		key := luaformatter.ValueToStringNoErr(a)
		e1, ok := e.entries[key]
		if !ok {
			e1 = newAssertEntry()
			e.entries[key] = e1
		}
		e = e1
	}

	e.asserts = append(e.asserts, called)

	return nil
}
