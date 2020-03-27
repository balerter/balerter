package mock

import (
	"fmt"
	"github.com/balerter/balerter/internal/lua_formatter"
	"github.com/balerter/balerter/internal/modules"
)

func (m *ModuleMock) Result() ([]modules.TestResult, error) {

	var result []modules.TestResult

	for _, e := range m.errors {
		result = append(result, modules.TestResult{
			ScriptName: "",
			ModuleName: m.name,
			Message:    e,
			Ok:         true,
		})
	}

	for hash, a := range m.assertsCalled {
		r := modules.TestResult{
			ScriptName: "",
			ModuleName: m.name,
			Message:    "assert called: ",
			Ok:         true,
		}

		calledCount, ok := m.queryLog[hash]

		if !ok {
			r.Message = "assert called fail: "
			r.Ok = false
		}
		if ok && calledCount != a.count {
			r.Message = "assert called count fail: "
			r.Ok = false
		}

		s, err := lua_formatter.ValuesToString(a.args)
		if err != nil {
			return nil, fmt.Errorf("error marshal lua.Value to a string, %w", err)
		}

		r.Message += fmt.Sprintf("count %d, expected %d, method '%s' with args %s", calledCount, a.count, a.method, s)

		result = append(result, r)
	}

	for hash, a := range m.assertsNotCalled {
		r := modules.TestResult{
			ScriptName: "",
			ModuleName: m.name,
			Message:    "assert not called: ",
			Ok:         true,
		}

		calledCount, ok := m.queryLog[hash]

		if ok {
			_ = a
			r.Message = "assert not called fail: "
			r.Ok = false
		}

		s, err := lua_formatter.ValuesToString(a.args)
		if err != nil {
			return nil, fmt.Errorf("error marshal lua.Value to a string, %w", err)
		}

		r.Message += fmt.Sprintf("count %d, expected %d, method '%s' with args %s", calledCount, a.count, a.method, s)

		result = append(result, r)
	}

	return result, nil
}
