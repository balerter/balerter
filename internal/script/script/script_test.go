package script

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScript_ParseMeta(t *testing.T) {
	s := &Script{
		Body: []byte(`
-- header
-- @interval 5m
-- foo
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
}

func TestScript_ParseMeta_WrongDuration(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @interval 5sm`),
	}

	err := s.ParseMeta()

	require.Error(t, err)
	assert.Equal(t, "time: unknown unit sm in duration 5sm", err.Error())
}
