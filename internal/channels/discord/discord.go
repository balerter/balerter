package discord

import (
	discordCfg "github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/session"
	"go.uber.org/zap"
)

type isession interface {
	SendMessage(channelID discord.Snowflake, content string, embed *discord.Embed) (*discord.Message, error)
}

// Discord implements a Provider for discord notifications.
type Discord struct {
	logger  *zap.Logger
	name    string
	session isession
	chanID  discord.Snowflake
}

// New returns the new Discord instance
func New(cfg discordCfg.Discord, logger *zap.Logger) (*Discord, error) {
	s, err := session.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	d := &Discord{
		logger:  logger,
		name:    cfg.Name,
		session: s,
		chanID:  discord.Snowflake(cfg.ChannelID),
	}

	return d, nil
}

// Name returns the Discord channel name
func (d *Discord) Name() string {
	return d.name
}
