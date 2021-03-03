package discord

import (
	"github.com/balerter/balerter/internal/config/channels/discord"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	d, err := New(discord.Discord{Name: "foo", ChannelID: 123}, nil)
	require.NoError(t, err)
	assert.IsType(t, &Discord{}, d)
	assert.Equal(t, "foo", d.name)
	assert.Equal(t, "123", d.chanID.String())
}

func TestName(t *testing.T) {
	d := &Discord{name: "foo"}
	assert.Equal(t, "foo", d.Name())
}
