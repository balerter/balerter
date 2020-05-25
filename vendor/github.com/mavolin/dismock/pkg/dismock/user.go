package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// User mocks a User request.
//
// The ID field of the passed User must be set.
func (m *Mocker) User(u discord.User) {
	m.MockAPI("User", http.MethodGet, "/users/"+u.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, u)
		})
}

// Me mocks a Me request.
//
// This method will sanitize User.ID.
func (m *Mocker) Me(u discord.User) {
	u = sanitize.User(u, 1)

	m.MockAPI("Me", http.MethodGet, "/users/@me",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, u)
		})
}

// ModifyMe mocks a ModifyMe request.
//
// This method will sanitize User.ID.
func (m *Mocker) ModifyMe(d api.ModifySelfData, u discord.User) {
	u = sanitize.User(u, 1)

	m.MockAPI("ModifyMe", http.MethodPatch, "/users/@me",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifySelfData), &d)
			mockutil.WriteJSON(t, w, u)
		})
}

type changeOwnNicknamePayload struct {
	Nick string `json:"nick"`
}

// ChangeOwnNickname mocks a ChangeOwnNickname request.
func (m *Mocker) ChangeOwnNickname(guildID discord.Snowflake, nick string) {
	m.MockAPI("ChangeOwnNickname", http.MethodPatch, "/guilds/"+guildID.String()+"/members/@me/nick",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := changeOwnNicknamePayload{
				Nick: nick,
			}

			mockutil.CheckJSON(t, r.Body, new(changeOwnNicknamePayload), &expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// PrivateChannels mocks a PrivateChannels request.
//
// This method will sanitize Channels.ID.
func (m *Mocker) PrivateChannels(c []discord.Channel) {
	for i, channel := range c {
		c[i] = sanitize.Channel(channel, 1)
	}

	m.MockAPI("PrivateChannels", http.MethodGet, "/users/@me/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}

type createPrivateChannelPayload struct {
	RecipientID discord.Snowflake `json:"recipient_id"`
}

// CreatePrivateChannel mocks a CreatePrivateChannel request.
//
// The c.DMRecipients[0] field of the passed Channel must be set.
func (m *Mocker) CreatePrivateChannel(c discord.Channel) {
	m.MockAPI("CreatePrivateChannel", http.MethodPost, "/users/@me/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := createPrivateChannelPayload{
				RecipientID: c.DMRecipients[0].ID,
			}

			mockutil.CheckJSON(t, r.Body, new(createPrivateChannelPayload), &expect)
			mockutil.WriteJSON(t, w, c)
		})
}

// UserConnections mocks a UserConnections request.
func (m *Mocker) UserConnections(c []discord.Connection) {
	m.MockAPI("UserConnections", http.MethodGet, "/users/@me/connections",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}
