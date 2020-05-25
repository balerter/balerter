package sanitize

import "github.com/diamondburned/arikawa/discord"

// Emoji sanitizes the passed emoji.
//
// This function will sanitize Emoji.ID and Emoji.User.ID.
func Emoji(e discord.Emoji, id, userID discord.Snowflake) discord.Emoji {
	if e.ID <= 0 {
		e.ID = id
	}

	e.User = User(e.User, userID)

	return e
}
