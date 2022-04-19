package dismock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// Invite mocks a Invite request.
//
// The Code field of the passed discord.Invite must be set.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func (m *Mocker) Invite(i discord.Invite) {
	i = sanitize.Invite(i, 1, 1, 1, 1, 1, 1, 1)

	m.MockAPI("Invite", http.MethodGet, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}

// InviteWithCounts mocks a InviteWithCounts request.
//
// The Code field of the passed discord.Invite must be set.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func (m *Mocker) InviteWithCounts(i discord.Invite) {
	i = sanitize.Invite(i, 1, 1, 1, 1, 1, 1, 1)

	m.MockAPI("Invite", http.MethodGet, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"with_counts": {"true"},
			})
			mockutil.WriteJSON(t, w, i)
		})
}

// ChannelInvites mocks a ChannelInvites request.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func (m *Mocker) ChannelInvites(channelID discord.ChannelID, invites []discord.Invite) {
	if invites == nil {
		invites = []discord.Invite{}
	}

	for i, invite := range invites {
		invites[i] = sanitize.Invite(invite, 1, 1, channelID, 1, 1, 1, 1)
	}

	m.MockAPI("ChannelInvites", http.MethodGet, "/channels/"+channelID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, invites)
		})
}

// GuildInvites mocks a GuildInvites request.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func (m *Mocker) GuildInvites(guildID discord.GuildID, invites []discord.Invite) {
	if invites == nil {
		invites = []discord.Invite{}
	}

	for i, invite := range invites {
		invites[i] = sanitize.Invite(invite, guildID, 1, 1, 1, 1, 1, 1)
	}

	m.MockAPI("GuildInvites", http.MethodGet, "/guilds/"+guildID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, invites)
		})
}

// CreateInvite mocks a CreateInvite request.
//
// The Channel.ID field of the passed discord.Invite must be set.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, and Invite.Inviter.ID.
func (m *Mocker) CreateInvite(d api.CreateInviteData, i discord.Invite) {
	i = sanitize.Invite(i, 1, 1, 1, 1, 1, 1, 1)

	m.MockAPI("CreateInvite", http.MethodPost, "/channels/"+i.Channel.ID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateInviteData), &d)
			mockutil.WriteJSON(t, w, i)
		})
}

// DeleteInvite mocks a DeleteInvite request.
//
// The Code field of the passed discord.Invite must be set.
//
// This method will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID and
// Invite.Inviter.ID.
func (m *Mocker) DeleteInvite(i discord.Invite) {
	i = sanitize.Invite(i, 1, 1, 1, 1, 1, 1, 1)

	m.MockAPI("DeleteInvite", http.MethodDelete, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}
