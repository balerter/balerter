package manager

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManager_Get(t *testing.T) {
	m := &Manager{storages: map[string]coreStorage.CoreStorage{
		"foo": nil,
	}}

	_, err := m.Get("foo")
	assert.NoError(t, err)

	_, err = m.Get("bar")
	require.Error(t, err)
	assert.Equal(t, "storage not found", err.Error())
}
