package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestRegistry_getAssetsOrphans_empty(t *testing.T) {
	r := New()
	results := r.getAssertsOrphans()
	assert.Equal(t, 0, len(results))
}

func TestRegistry_getAssetsOrphans_NoArgs(t *testing.T) {
	r := New()
	e := newAssertEntry()
	e.asserts = append(e.asserts, true, false)
	r.assertEntries["foo"] = e

	results := r.getAssertsOrphans()
	assert.Equal(t, 2, len(results))

	assert.Equal(t, false, results[0].Ok)
	assert.Equal(t, "method 'foo' with args [] was not called, but should", results[0].Message)

	assert.Equal(t, true, results[1].Ok)
	assert.Equal(t, "method 'foo' with args [] was not called", results[1].Message)
}

func TestRegistry_getAssetsOrphans_WithArgs(t *testing.T) {
	r := New()
	e2 := newAssertEntry()
	e2.asserts = append(e2.asserts, true, false)

	e1 := newAssertEntry()
	e1.asserts = append(e1.asserts, true, false)
	e1.entries["baz"] = e2

	e := newAssertEntry()
	e.asserts = append(e.asserts, true, false)
	e.entries["bar"] = e1

	r.assertEntries["foo"] = e

	results := r.getAssertsOrphans()
	require.Equal(t, 6, len(results))

	assert.Equal(t, false, results[0].Ok)
	assert.Equal(t, "method 'foo' with args [] was not called, but should", results[0].Message)

	assert.Equal(t, true, results[1].Ok)
	assert.Equal(t, "method 'foo' with args [] was not called", results[1].Message)

	assert.Equal(t, false, results[2].Ok)
	assert.Equal(t, "method 'foo' with args [bar] was not called, but should", results[2].Message)

	assert.Equal(t, true, results[3].Ok)
	assert.Equal(t, "method 'foo' with args [bar] was not called", results[3].Message)

	assert.Equal(t, false, results[4].Ok)
	assert.Equal(t, "method 'foo' with args [bar baz] was not called, but should", results[4].Message)

	assert.Equal(t, true, results[5].Ok)
	assert.Equal(t, "method 'foo' with args [bar baz] was not called", results[5].Message)
}

func TestRegistry_getAssert(t *testing.T) {
	r := New()

	e, ok := r.getAssert(call{"foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}})
	assert.False(t, ok)
	assert.Nil(t, e)

	ent := newAssertEntry()
	r.assertEntries["foo"] = ent

	e, ok = r.getAssert(call{"foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}})
	assert.False(t, ok)
	assert.Nil(t, e)

	ent1 := newAssertEntry()
	ent.entries["bar"] = ent1

	e, ok = r.getAssert(call{"foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}})
	assert.False(t, ok)
	assert.Nil(t, e)

	ent2 := newAssertEntry()
	ent1.entries["42"] = ent2

	e, ok = r.getAssert(call{"foo", []lua.LValue{lua.LString("bar"), lua.LNumber(42)}})
	assert.True(t, ok)
	assert.Equal(t, ent2, e)
}

func TestResult(t *testing.T) {
	r := New()

	r.calls = append(r.calls,
		call{method: "f1", args: []lua.LValue{lua.LString("foo")}}, // no asserts
		call{method: "f11", args: nil},                             // no asserts
		call{method: "f2", args: []lua.LValue{lua.LString("foo")}}, // assert called
		call{method: "f3", args: []lua.LValue{lua.LString("foo")}}, // assert not called
	)

	// assert called for 'f2'
	e2 := newAssertEntry()
	e2.asserts = append(e2.asserts, true)
	e1 := newAssertEntry()
	e1.entries["foo"] = e2
	r.assertEntries["f2"] = e1

	// assert not called for 'f3'
	e2 = newAssertEntry()
	e2.asserts = append(e2.asserts, false)
	e1 = newAssertEntry()
	e1.entries["foo"] = e2
	r.assertEntries["f3"] = e1

	e1 = newAssertEntry()
	r.assertEntries["f11"] = e1

	results := r.Result()
	require.Equal(t, 2, len(results))

	assert.Equal(t, true, results[0].Ok)
	assert.Equal(t, "method 'f2' with args [foo] was called", results[0].Message)

	assert.Equal(t, false, results[1].Ok)
	assert.Equal(t, "method 'f3' with args [foo] was called, but should not", results[1].Message)
}
