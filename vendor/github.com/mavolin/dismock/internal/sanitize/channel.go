package sanitize

import "github.com/diamondburned/arikawa/discord"

// Channel sanitizes a discord.Channel.
//
// This function will sanitize Channel.ID.
func Channel(c discord.Channel, id discord.ChannelID) discord.Channel {
	if c.ID == 0 {
		c.ID = id
	}

	return c
}

// Message sanitizes a discord.Message.
//
// This function will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func Message(
	m discord.Message, id discord.MessageID, channelID discord.ChannelID, authorID discord.UserID,
) discord.Message {
	if m.ID == 0 {
		m.ID = id
	}

	if m.ChannelID == 0 {
		m.ChannelID = channelID
	}

	m.Author = User(m.Author, authorID)

	return m
}
