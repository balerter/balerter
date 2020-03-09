package provider

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemory_Put(t *testing.T) {
	m := New()

	err := m.Put("foo", "bar")
	require.NoError(t, err)

	err = m.Put("foo", "bar2")
	require.Error(t, err)

	v, ok := m.storage["foo"]
	assert.True(t, ok)
	assert.Equal(t, "bar", v)
}

func TestMemory_Upsert(t *testing.T) {
	m := New()

	err := m.Upsert("foo", "bar")
	require.NoError(t, err)

	err = m.Upsert("foo", "bar2")
	require.NoError(t, err)

	v, ok := m.storage["foo"]
	assert.True(t, ok)
	assert.Equal(t, "bar2", v)
}

func TestMemory_Delete(t *testing.T) {
	m := New()

	m.storage["foo"] = "bar"

	err := m.Delete("foo2")
	require.Error(t, err)

	err = m.Delete("foo")
	require.NoError(t, err)

	_, ok := m.storage["foo"]
	assert.False(t, ok)

	err = m.Delete("foo")
	require.Error(t, err)
}

func TestMemory_Get(t *testing.T) {
	m := New()

	v, err := m.Get("foo")
	require.Error(t, err)

	m.storage["foo"] = "bar"

	v, err = m.Get("foo")
	require.NoError(t, err)

	assert.Equal(t, "bar", v)
}
