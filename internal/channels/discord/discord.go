package discord

import (
	discordCfg "github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/session"
	"go.uber.org/zap"
)

//go:generate moq -out module_mock_session.go -skip-ensure -fmt goimports . isession

type isession interface {
	SendMessage(channelID discord.ChannelID, content string, embed *discord.Embed) (*discord.Message, error)
}

// Discord implements a Provider for discord notifications.
type Discord struct {
	logger  *zap.Logger
	name    string
	session isession
	chanID  discord.ChannelID
	ignore  bool
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
		chanID:  discord.ChannelID(cfg.ChannelID),
		ignore:  cfg.Ignore,
	}

	return d, nil
}

// Name returns the Discord channel name
func (d *Discord) Name() string {
	return d.name
}

func (d *Discord) Ignore() bool {
	return d.ignore
}
