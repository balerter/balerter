package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_registry_response_simple(t *testing.T) {
	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("foo")}

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			responses: [][]lua.LValue{ret1},
		},
	}

	resp, err := reg.Response("AnyValue", "foo", nil)
	require.NoError(t, err)
	assert.Equal(t, ret1, resp)
}

func Test_registry_response_tree1(t *testing.T) {
	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("foo")}

	arg1 := lua.LNumber(42)
	arg2 := lua.LString("bar")
	arg3 := &lua.LTable{}
	arg3.RawSetString("foo", lua.LNumber(10))

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			entries: map[string]*responseEntry{
				"number@42": {
					entries: map[string]*responseEntry{
						"string@bar": {
							entries: map[string]*responseEntry{
								"table@{\"foo\":10}": {
									responses: [][]lua.LValue{ret1},
								},
							},
						},
					},
				},
			},
		},
	}

	resp, err := reg.Response("AnyValue", "foo", []lua.LValue{arg1, arg2, arg3})
	require.NoError(t, err)
	assert.Equal(t, ret1, resp)
}

func Test_registry_response_tree1_two_values(t *testing.T) {
	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("foo")}
	ret2 := []lua.LValue{lua.LNumber(58), lua.LString("bar")}

	arg1 := lua.LNumber(42)
	arg2 := lua.LString("bar")
	arg3 := &lua.LTable{}
	arg3.RawSetString("foo", lua.LNumber(10))

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			entries: map[string]*responseEntry{
				"number@42": {
					entries: map[string]*responseEntry{
						"string@bar": {
							entries: map[string]*responseEntry{
								"table@{\"foo\":10}": {
									responses: [][]lua.LValue{ret2, ret1},
								},
							},
						},
					},
				},
			},
		},
	}

	resp, err := reg.Response("AnyValue", "foo", []lua.LValue{arg1, arg2, arg3})
	require.NoError(t, err)
	assert.Equal(t, ret1, resp)

	resp, err = reg.Response("AnyValue", "foo", []lua.LValue{arg1, arg2, arg3})
	require.NoError(t, err)
	assert.Equal(t, ret2, resp)
}

func Test_registry_response_any_value(t *testing.T) {
	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("foo")}

	arg1 := lua.LNumber(42)
	arg2 := lua.LString("bar")
	arg3 := &lua.LTable{}
	arg3.RawSetString("foo", lua.LNumber(10))

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			entries: map[string]*responseEntry{
				"number@42": {
					entries: map[string]*responseEntry{
						"AnyValue": {
							entries: map[string]*responseEntry{
								"table@{\"foo\":10}": {
									responses: [][]lua.LValue{ret1},
								},
							},
						},
					},
				},
			},
		},
	}

	resp, err := reg.Response("AnyValue", "foo", []lua.LValue{arg1, arg2, arg3})
	require.NoError(t, err)
	assert.Equal(t, ret1, resp)
}

func Test_registry_response_error1(t *testing.T) {
	// save one response but require two times
	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("foo")}

	arg1 := lua.LNumber(42)

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			entries: map[string]*responseEntry{
				"number@42": {
					responses: [][]lua.LValue{
						ret1,
					},
				},
			},
		},
	}

	resp, err := reg.Response("AnyValue", "foo", []lua.LValue{arg1})
	require.NoError(t, err)
	assert.Equal(t, ret1, resp)

	resp, err = reg.Response("AnyValue", "foo", []lua.LValue{arg1})
	require.Error(t, err)
}

func Test_registry_response_no_entry_in_tree(t *testing.T) {
	arg1 := lua.LNumber(43)

	reg := New()
	reg.responseEntries = map[string]*responseEntry{
		"foo": {
			entries: map[string]*responseEntry{
				"number@42": {
					responses: [][]lua.LValue{},
				},
			},
		},
	}

	_, err := reg.Response("AnyValue", "foo", []lua.LValue{arg1})
	require.Error(t, err)
}

func Test_registry_response_error_no_method(t *testing.T) {
	reg := New()

	_, err := reg.Response("AnyValue", "foo", []lua.LValue{})
	require.Error(t, err)
}
