package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistryClean(t *testing.T) {
	r := New()
	r.calls = append(r.calls, call{}, call{})
	r.assertEntries["foo"] = newAssertEntry()
	r.responseEntries["foo"] = newResponseEntry()

	r.Clean()

	assert.Equal(t, 0, len(r.calls))
	assert.Equal(t, 0, len(r.assertEntries))
	assert.Equal(t, 0, len(r.responseEntries))
}
