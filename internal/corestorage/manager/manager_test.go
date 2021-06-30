package manager

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func TestManager_Get_empty(t *testing.T) {
	m := &Manager{storages: map[string]coreStorage.CoreStorage{
		"memory": nil,
	}}

	_, err := m.Get("")
	assert.NoError(t, err)
}

func TestManager_Stop(t *testing.T) {
	s := &coreStorage.Mock{}
	s.On("Stop").Return(nil)

	m := &Manager{
		storages: map[string]coreStorage.CoreStorage{
			"foo": s,
		},
		logger: zap.NewNop(),
	}

	m.Stop()

	s.AssertCalled(t, "Stop")
}
