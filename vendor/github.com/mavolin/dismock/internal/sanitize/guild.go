package sanitize

import "github.com/diamondburned/arikawa/discord"

// Guild sanitizes a Guild.
//
// This function will sanitize Guild.ID, Guild.OwnerID, Guild.RulesChannelID,
// Guild.PublicUpdatesChannelID, Guild.Emojis.ID and Guild.Roles.ID.
func Guild(g discord.Guild, id, ownerID, rulesChannelID, publicUpdatesChannelID discord.Snowflake) discord.Guild {
	if g.ID <= 0 {
		g.ID = id
	}

	if g.OwnerID <= 0 {
		g.OwnerID = ownerID
	}

	if g.RulesChannelID <= 0 {
		g.RulesChannelID = rulesChannelID
	}

	if g.PublicUpdatesChannelID <= 0 {
		g.PublicUpdatesChannelID = publicUpdatesChannelID
	}

	for i, emoji := range g.Emojis {
		g.Emojis[i] = Emoji(emoji, 1, 1)
	}

	for i, role := range g.Roles {
		g.Roles[i] = Role(role, 1)
	}

	return g
}

// GuildPreview sanitizes a GuildPreview.
//
// This function will sanitize GuildPreview.ID, GuildPreview.Emojis.ID and
// GuildPreview.Emojis.User.ID.
func GuildPreview(p discord.GuildPreview, id discord.Snowflake) discord.GuildPreview {
	if p.ID <= 0 {
		p.ID = id
	}

	for i, emoji := range p.Emojis {
		p.Emojis[i] = Emoji(emoji, 1, 1)
	}

	return p
}

// Role sanitizes a Role.
//
// This function will sanitize Role.ID.
func Role(r discord.Role, id discord.Snowflake) discord.Role {
	if r.ID <= 0 {
		r.ID = id
	}

	return r
}

// Member sanitizes a Member.
//
// This function will sanitize Member.User.ID.
func Member(m discord.Member, userID discord.Snowflake) discord.Member {
	if m.User.ID <= 0 {
		m.User.ID = userID
	}

	return m
}

// Ban sanitizes a Ban.
//
// This function will sanitize Ban.User.ID.
func Ban(b discord.Ban, userID discord.Snowflake) discord.Ban {
	b.User = User(b.User, userID)

	return b
}

// Integration sanitizes an Integration.
//
// This function will sanitize Integration.ID, Integration.RoleID and
// Integration.User.ID.
func Integration(i discord.Integration, id, roleID, userID discord.Snowflake) discord.Integration {
	if i.ID <= 0 {
		i.ID = id
	}

	if i.RoleID <= 0 {
		i.RoleID = roleID
	}

	i.User = User(i.User, userID)

	return i
}
