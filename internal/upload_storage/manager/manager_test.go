package manager

import (
	"testing"

	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/balerter/balerter/internal/modules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	m := New(zap.NewNop())
	assert.IsType(t, &Manager{}, m)
}

func TestManager_Get(t *testing.T) {
	m1 := &modules.ModuleMock{NameFunc: func() string {
		return "foo"
	}}
	m2 := &modules.ModuleMock{NameFunc: func() string {
		return "bar"
	}}

	m := &Manager{
		logger:  zap.NewNop(),
		modules: map[string]modules.Module{"foo": m1, "bar": m2},
	}

	result := m.Get()
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, m1)
	assert.Contains(t, result, m2)
}

func TestManager_Init(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.Module{},
	}

	err := m.Init(&upload.Upload{S3: []s3.S3{{
		Name:     "f1",
		Region:   "f2",
		Key:      "f3",
		Secret:   "f4",
		Endpoint: "f5",
		Bucket:   "f6",
	}}})

	require.NoError(t, err)

	md, ok := m.modules["s3.f1"]
	require.True(t, ok)

	assert.Equal(t, "s3.f1", md.Name())
}
