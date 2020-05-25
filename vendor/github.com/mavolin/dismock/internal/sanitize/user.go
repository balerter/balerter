package sanitize

import "github.com/diamondburned/arikawa/discord"

// User sanitizes a User.
//
// This function will sanitize User.ID.
func User(u discord.User, id discord.Snowflake) discord.User {
	if u.ID <= 0 {
		u.ID = id
	}

	return u
}
