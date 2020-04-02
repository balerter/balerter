package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestRegistry_AddAssert1(t *testing.T) {
	r := New()
	err := r.AddAssert("foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}, true)
	require.NoError(t, err)
	err = r.AddAssert("foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}, false)
	require.NoError(t, err)

	assert.Equal(t, 1, len(r.assertEntries))
	e, ok := r.assertEntries["foo"]
	require.True(t, ok)
	assert.Equal(t, 0, len(e.asserts))
	assert.Equal(t, 1, len(e.entries))
	e1, ok := e.entries["bar"]
	require.True(t, ok)
	assert.Equal(t, 0, len(e1.asserts))
	assert.Equal(t, 1, len(e1.entries))
	e2, ok := e1.entries["42"]
	require.True(t, ok)
	assert.Equal(t, 2, len(e2.asserts))
	assert.Equal(t, 0, len(e2.entries))
	assert.Equal(t, true, e2.asserts[0])
	assert.Equal(t, false, e2.asserts[1])
}
