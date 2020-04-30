package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScript_ParseMeta_Without_Meta(t *testing.T) {
	everySecondSched, err := NewSchedule("@every 1s")
	require.NoError(t, err)
	s := &Script{
		Schedule: everySecondSched,
		Body: []byte(`
print
-- @schedule @every 1s
-- @ignore
`),
	}

	err = s.parseMeta()
	require.NoError(t, err)

	assert.False(t, s.Ignore)
	assert.Equal(t, everySecondSched, s.Schedule)

	s = &Script{
		Schedule: everySecondSched,
		Body: []byte(`
-- hello 
`),
	}

	err = s.parseMeta()

	require.NoError(t, err)

	assert.False(t, s.Ignore)
	assert.Equal(t, everySecondSched, s.Schedule)
}

func TestScript_ParseMeta_Multiple_Crons(t *testing.T) {
	s := &Script{
		Body: []byte(`
-- header
-- @every 5m
-- foo
-- @name newname
--
-- @ignore
-- @every 6m
--
print
-- @every 10m
`),
	}

	err := s.parseMeta()
	cronErr, ok := err.(CronAlreadyDefinedError)
	if !ok {
		t.Error("unexpected error:", err)
	}

	assert.Equal(t, cronErr.prevSpec, "@every 5m")
}

func TestScript_ParseMeta_WrongDuration(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @every 5sm`),
	}

	err := s.parseMeta()

	require.Error(t, err)
	assert.Equal(t, "failed to parse duration @every 5sm: time: unknown unit sm in duration 5sm", err.Error())
}

func TestScript_ParseMeta_EmptyName(t *testing.T) {
	s := &Script{
		Body: []byte(`-- @name  `),
	}

	err := s.parseMeta()

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

	err := s.parseMeta()

	require.Error(t, err)
	assert.Equal(t, "channels must be not empty", err.Error())
}

func TestScript_ParseMeta_EmptyChannel(t *testing.T) {
	s := &Script{Body: []byte(`-- @channels foo,`)}
	err := s.parseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())

	s = &Script{Body: []byte(`-- @channels ,foo`)}
	err = s.parseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())

	s = &Script{Body: []byte(`-- @channels bar,,foo`)}
	err = s.parseMeta()
	require.Error(t, err)
	assert.Equal(t, "channel name must be not empty", err.Error())
}

func TestScript_ParseMeta_Channels(t *testing.T) {
	s := &Script{Body: []byte(`-- @channels foo`)}
	err := s.parseMeta()
	require.NoError(t, err)
	assert.Equal(t, 1, len(s.Channels))
	assert.Contains(t, s.Channels, "foo")

	s = &Script{Body: []byte(`-- @channels foo, bar`)}
	err = s.parseMeta()
	require.NoError(t, err)
	assert.Equal(t, 2, len(s.Channels))
	assert.Contains(t, s.Channels, "foo")
	assert.Contains(t, s.Channels, "bar")
}
