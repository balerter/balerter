package test

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
)

func (m *Mock) Result() []modules.TestResult {
	var result []modules.TestResult

	for _, e := range m.errors {
		result = append(result, modules.TestResult{
			ModuleName: m.name,
			Message:    e,
			Ok:         false,
		})
	}

	for q, cnt := range m.asserts.callQuery {
		c, ok := m.queryLog[q]
		message := fmt.Sprintf("assertCalledQuery: the query '%s' has been called %d times", q, c)
		pass := true

		if !ok {
			message = fmt.Sprintf("assertCalledQuery: the query '%s' has not been called", q)
			pass = false
		}

		if ok && c < cnt {
			message = fmt.Sprintf("assertCalledQuery: the query '%s' has been called %d times, expected %d", q, c, cnt)
			pass = false
		}

		result = append(result, modules.TestResult{
			ModuleName: m.name,
			Message:    message,
			Ok:         pass,
		})
	}

	return result
}
