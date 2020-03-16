package memory

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlert_GetOrNew(t *testing.T) {
	m := New()

	a, err := m.GetOrNew("foo")
	assert.NoError(t, err)
	assert.IsType(t, &alert.Alert{}, a)

	assert.Equal(t, 1, len(m.alerts))
	_, ok := m.alerts["foo"]
	assert.True(t, ok)

	a2, err := m.GetOrNew("foo")
	assert.NoError(t, err)
	assert.IsType(t, &alert.Alert{}, a2)

	assert.Equal(t, 1, len(m.alerts))
	_, ok = m.alerts["foo"]
	assert.True(t, ok)
	assert.Equal(t, a, a2)
}

func TestAlert_All(t *testing.T) {
	m := New()

	a1 := &alert.Alert{}
	a1.SetName("foo")
	m.alerts["foo"] = a1

	a2 := &alert.Alert{}
	a2.SetName("bar")
	m.alerts["bar"] = a2

	aa := m.All()
	require.Equal(t, 2, len(aa))
	assert.Contains(t, aa, a1)
	assert.Contains(t, aa, a2)
}
