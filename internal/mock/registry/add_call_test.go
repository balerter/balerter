package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestRegistry_AddCall(t *testing.T) {
	r := New()

	args1 := []lua.LValue{lua.LString("bar1"), lua.LNumber(42)}
	args2 := []lua.LValue{lua.LString("bar2"), lua.LNumber(43)}

	err := r.AddCall("foo1", args1)
	require.NoError(t, err)
	err = r.AddCall("foo2", args2)
	require.NoError(t, err)

	assert.Equal(t, 2, len(r.calls))
	assert.Equal(t, "foo1", r.calls[0].method)
	assert.Equal(t, args1, r.calls[0].args)
	assert.Equal(t, "foo2", r.calls[1].method)
	assert.Equal(t, args2, r.calls[1].args)
}
