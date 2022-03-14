package sanitize

import "github.com/diamondburned/arikawa/discord"

// Emoji sanitizes an discord.Emoji.
//
// This function will sanitize Emoji.ID and Emoji.User.ID.
func Emoji(e discord.Emoji, id discord.EmojiID, userID discord.UserID) discord.Emoji {
	if e.ID == 0 {
		e.ID = id
	}

	e.User = User(e.User, userID)

	return e
}
