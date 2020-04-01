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

	for _, res := range m.registry.Result() {
		res.ModuleName = m.name
		result = append(result, res)
	}

	return result, nil
}
