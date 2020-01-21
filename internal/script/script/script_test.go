package script

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScript_ParseMeta_without_meta(t *testing.T) {
	s := &Script{
		Interval: time.Second,
		Body: []byte(`
print
-- @interval 10m
-- @ignore
`),
	}

	err := s.ParseMeta()

	require.NoError(t, err)

	assert.False(t, s.Ignore)
	assert.Equal(t, time.Second, s.Interval)
}

func TestScript_ParseMeta(t *testing.T) {
	s := &Script{
		Body: []byte(`
-- header
-- @interval 5m
-- foo
-- @name newname
--
-- @ignore
-- @interval 6m
--
print
-- @interval 10m
`),
	}

	err := s.ParseMeta()

	require.NoError(t, err)

	assert.True(t, s.Ignore)
	assert.Equal(t, time.Minute*6, s.Interval)
	assert.Equal(t, "newname", s.Name)
}

func TestScript_ParseMeta_WrongDuration(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @interval 5sm`),
	}

	err := s.ParseMeta()

	require.Error(t, err)
	assert.Equal(t, "time: unknown unit sm in duration 5sm", err.Error())
}

func TestScript_ParseMeta_EmptyName(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @name  `),
	}

	err := s.ParseMeta()

	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())
}

func TestScript_Hash(t *testing.T) {
	s := &Script{
		Name: "name",
		Body: []byte("body"),
	}

	assert.Equal(t, "e1073eb510ec5486c4f5a62a1ba486202db8aa25", s.Hash())
}
