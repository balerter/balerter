package script

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScript_error_parse_cron_value(t *testing.T) {
	s := &Script{
		Body: []byte(`
-- @cron 0 0 0 0 0 0
`),
	}

	err := s.ParseMeta()

	require.Error(t, err)
	assert.Equal(t, "error parse cron value, beginning of range (0) below minimum (1): 0", err.Error())
}

func TestScript_ParseMeta_without_meta(t *testing.T) {
	s := &Script{
		CronValue: "* * * * * *",
		Body: []byte(`
print
-- @cron 1 1 1 1 1 1
-- @ignore
`),
	}

	err := s.ParseMeta()

	require.NoError(t, err)

	assert.False(t, s.Ignore)
	assert.Equal(t, "* * * * * *", s.CronValue)

	s = &Script{
		CronValue: "* * * * * *",
		Body: []byte(`
-- hello 
`),
	}

	err = s.ParseMeta()

	require.NoError(t, err)

	assert.False(t, s.Ignore)
	assert.Equal(t, "* * * * * *", s.CronValue)
}

func TestScript_ParseMeta(t *testing.T) {
	s := &Script{
		Body: []byte(`
-- header
-- @cron 1 1 1 1 1 1
-- foo
-- @name newname
--
-- @ignore
-- @interval 6m
-- @timeout 10s
--
print
-- @interval 10m
`),
	}

	err := s.ParseMeta()

	require.NoError(t, err)

	assert.True(t, s.Ignore)
	assert.Equal(t, "1 1 1 1 1 1", s.CronValue)
	assert.Equal(t, "newname", s.Name)
	assert.Equal(t, time.Second*10, s.Timeout)
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

func TestScript_ParseMeta_EmptyChannels(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @channels `),
	}

	err := s.ParseMeta()

	require.Error(t, err)
	assert.Equal(t, "channels must be not empty", err.Error())
}

func TestScript_ParseMeta_EmptyChannel(t *testing.T) {
	s := &Script{Body: []byte(`-- @channels foo,`)}
	err := s.ParseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())

	s = &Script{Body: []byte(`-- @channels ,foo`)}
	err = s.ParseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())

	s = &Script{Body: []byte(`-- @channels bar,,foo`)}
	err = s.ParseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())
}

func TestScript_ParseMeta_Channels(t *testing.T) {
	s := &Script{Body: []byte(`-- @channels foo`)}
	err := s.ParseMeta()
	require.NoError(t, err)
	assert.Equal(t, 1, len(s.Channels))
	assert.Contains(t, s.Channels, "foo")

	s = &Script{Body: []byte(`-- @channels foo, bar`)}
	err = s.ParseMeta()
	require.NoError(t, err)
	assert.Equal(t, 2, len(s.Channels))
	assert.Contains(t, s.Channels, "foo")
	assert.Contains(t, s.Channels, "bar")
}
