package sanitize

import "github.com/diamondburned/arikawa/discord"

// Channel sanitizes a Channel.
//
// This function will sanitize Channel.ID.
func Channel(c discord.Channel, id discord.Snowflake) discord.Channel {
	if c.ID <= 0 {
		c.ID = id
	}

	return c
}

// Message sanitizes a message.
//
// This function will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func Message(m discord.Message, id, channelID, authorID discord.Snowflake) discord.Message {
	if m.ID <= 0 {
		m.ID = id
	}

	if m.ChannelID <= 0 {
		m.ChannelID = channelID
	}

	m.Author = User(m.Author, authorID)

	return m
}
