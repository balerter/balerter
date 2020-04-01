package registry

import (
	"fmt"
	"github.com/balerter/balerter/internal/lua_formatter"
	"github.com/balerter/balerter/internal/modules"
)

func (r *Registry) Result() []modules.TestResult {
	results := make([]modules.TestResult, 0)

	// iterate over calls and pick asserts
	for _, call := range r.calls {
		// an assert for this key is not registered
		e, ok := r.getAssert(call)
		if !ok {
			continue
		}
		// all asserts already processed
		if len(e.asserts) == 0 {
			continue
		}

		assertShouldBeCalled := e.asserts[0]
		e.asserts = e.asserts[1:]

		if !assertShouldBeCalled {
			res := modules.TestResult{
				Message: fmt.Sprintf("method '%s' with args %s was called, but should not", call.method, lua_formatter.ValuesToStringNoErr(call.args)),
				Ok:      false,
			}
			results = append(results, res)
		}
	}

	for _, a := range r.getAssertsOrphans() {
		res := modules.TestResult{
			Message: fmt.Sprintf("method '%s' with args %v was not called, but should", a[0], a[1:]),
			Ok:      false,
		}
		results = append(results, res)
	}

	return results
}

func (r *Registry) getAssert(call call) (*assertEntry, bool) {
	e, ok := r.assertEntries[call.method]
	if !ok {
		return nil, false
	}

	for _, a := range call.args {
		key := lua_formatter.ValueToStringNoErr(a)
		e1, ok := e.entries[key]
		if !ok {
			return nil, false
		}
		e = e1
	}

	return e, true
}

// first element is a method name, next elements is args
func (r *Registry) getAssertsOrphans() [][]string {
	result := make([][]string, 0)

	for method, e := range r.assertEntries {
		for _, res := range r.getAssertsOrphansChain(e) {
			result = append(result, append([]string{method}, res...))
		}
	}

	return result
}

func (r *Registry) getAssertsOrphansChain(entry *assertEntry, args ...string) [][]string {
	result := make([][]string, 0)

	if len(entry.asserts) > 0 {
		for _, v := range entry.asserts {
			if v {
				result = append(result, args)
			}
		}
	}

	for arg, e := range entry.entries {
		res := r.getAssertsOrphansChain(e, append(args, arg)...)
		result = append(result, res...)
	}

	return result
}
