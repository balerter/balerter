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

	result = append(result, m.registry.Result()...)

	return result, nil
}
