package discord

import (
	discord2 "github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/session"
	"go.uber.org/zap"
)

// Discord implements a Provider for discord notifications.
type Discord struct {
	logger  *zap.Logger
	name    string
	session *session.Session
	chanID  discord.Snowflake
}

// New returns the new Discord instance
func New(cfg *discord2.Discord, logger *zap.Logger) (*Discord, error) {
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
