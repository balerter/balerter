package sanitize

import "github.com/diamondburned/arikawa/discord"

// Invite sanitizes an discord.Invite.
//
// This function will sanitizeInvite. Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.RulesChannelID, Invite.Guild.PublicUpdatesChannelID,
// Invite.Guild.Emojis.ID and Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func Invite(
	i discord.Invite, guildID discord.GuildID, ownerID discord.UserID,
	rulesChannelID, publicUpdatesChannelID, channelID discord.ChannelID, inviterID discord.UserID,
	targetID discord.UserID,
) discord.Invite {
	if i.Guild != nil {
		*i.Guild = Guild(*i.Guild, guildID, ownerID, rulesChannelID, publicUpdatesChannelID)
	}

	i.Channel = Channel(i.Channel, channelID)

	if i.Inviter != nil {
		*i.Inviter = User(*i.Inviter, inviterID)
	}

	if i.Target != nil {
		*i.Target = User(*i.Target, targetID)
	}

	return i
}
