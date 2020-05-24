package discord

import (
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/bwmarrin/discordgo"
)

// Send implements
func (d *Discord) Send(mes *message.Message) error {

	c, err := discordgo.New("Bot " + d.conf.Token)
	if err != nil {
		return err
	}

	_, err = c.ChannelMessageSend(d.conf.ChannelID, mes.Text)
	if err != nil {
		return err
	}

	return nil
}
