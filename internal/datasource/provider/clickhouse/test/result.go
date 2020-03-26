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
		message := fmt.Sprintf("assertCalledQuery: the query has been called %d times: %s", c, q)
		pass := true

		if !ok {
			message = fmt.Sprintf("assertCalledQuery: the query has not been called: %s", q)
			pass = false
		}

		if ok && c < cnt {
			message = fmt.Sprintf("assertCalledQuery: the query has been called %d times, expected %d: %s", c, cnt, q)
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
