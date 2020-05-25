package discord

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	tk := os.Getenv("DG_TOKEN")
	ci := os.Getenv("DG_CHANNEL_ID")

	if len(tk) == 0 || len(ci) == 0 {
		return
	}

	dg, err := discordgo.New("Bot " + tk)
	require.NoError(t, err)

	_, err = dg.ChannelMessageSend(ci, "Ping!")
	require.NoError(t, err)

}
