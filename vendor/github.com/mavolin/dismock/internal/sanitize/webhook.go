package sanitize

import "github.com/diamondburned/arikawa/discord"

// Webhook sanitizes a discord.Webhook.
//
// This function will sanitize Webhook.ID, Webhook.User.ID and
// Webhook.ChannelID.
func Webhook(
	w discord.Webhook, id discord.WebhookID, userID discord.UserID, channelID discord.ChannelID,
) discord.Webhook {
	if w.ID == 0 {
		w.ID = id
	}

	w.User = User(w.User, userID)

	if w.ChannelID == 0 {
		w.ChannelID = channelID
	}

	return w
}
