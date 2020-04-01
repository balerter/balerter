package mock

import (
	"github.com/balerter/balerter/internal/modules"
)

func (m *ModuleMock) Result() ([]modules.TestResult, error) {

	var result []modules.TestResult

	for _, e := range m.errors {
		result = append(result, modules.TestResult{
			ModuleName: m.name,
			Message:    e,
			Ok:         false,
		})
	}

	// todo: fetch asserts errors from the registry
	//for _, res := range m.queryLog.Asserts {
	//
	//	mes := fmt.Sprintf(
	//		"assert: method '%s' with args '%s'. want called: %t, assert %d times, called %d times",
	//		res.Method,
	//		lua_formatter.ValuesToStringNoErr(res.Args),
	//		res.WantCalled,
	//		res.AssertsCount,
	//		res.CallsCount,
	//	)
	//
	//	r := modules.TestResult{
	//		ModuleName: m.name,
	//		Message:    mes,
	//	}
	//
	//	//if res.WantCalled && res.
	//
	//	result = append(result, r)
	//}

	return result, nil
}
