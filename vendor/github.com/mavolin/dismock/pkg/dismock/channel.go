package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// Channels mocks a channels request.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channels(guildID discord.Snowflake, c []discord.Channel) {
	if c == nil {
		c = []discord.Channel{}
	}

	for i, channel := range c {
		c[i] = sanitize.Channel(channel, 1)
	}

	m.MockAPI("Channels", http.MethodGet, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}

// CreateChannel mocks a CreateChannel request.
//
// The GuildID field of the passed Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) CreateChannel(d api.CreateChannelData, c discord.Channel) {
	c = sanitize.Channel(c, 1)

	m.MockAPI("CreateChannel", http.MethodPost, "/guilds/"+c.GuildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateChannelData), &d)
			mockutil.WriteJSON(t, w, c)
		})
}

// MoveChannel mocks a MoveChannel request.
func (m *Mocker) MoveChannel(guildID discord.Snowflake, d []api.MoveChannelData) {
	m.MockAPI("CreateChannel", http.MethodPatch, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &[]api.MoveChannelData{}, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Channel mocks a Channel request.
//
// The ID field of the passed Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channel(c discord.Channel) {
	c = sanitize.Channel(c, 1)

	m.MockAPI("CreateChannel", http.MethodGet, "/channels/"+c.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}

// ModifyChannel mocks a ModifyChannel request.
func (m *Mocker) ModifyChannel(id discord.Snowflake, d api.ModifyChannelData) {
	m.MockAPI("ModifyChannel", http.MethodPatch, "/channels/"+id.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyChannelData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannel mocks a DeleteChannel request.
func (m *Mocker) DeleteChannel(id discord.Snowflake) {
	m.MockAPI("DeleteChannel", http.MethodDelete, "/channels/"+id.String(), nil)
}

// EditChannelPermission mocks a EditChannelPermission request.
func (m *Mocker) EditChannelPermission(channelID, overwriteID discord.Snowflake, d api.EditChannelPermissionData) {
	m.MockAPI("EditChannelPermission", http.MethodPut,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.EditChannelPermissionData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannelPermission mocks a DeleteChannelPermission request.
func (m *Mocker) DeleteChannelPermission(channelID, overwriteID discord.Snowflake) {
	m.MockAPI("DeleteChannelPermission", http.MethodDelete,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(), nil)
}

// Typing mocks a Typing request.
func (m *Mocker) Typing(channelID discord.Snowflake) {
	m.MockAPI("Typing", http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
}

// PinnedMessages mocks a PinnedMessages request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) PinnedMessages(channelID discord.Snowflake, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	for i, message := range messages {
		messages[i] = sanitize.Message(message, 1, channelID, 1)
	}

	m.MockAPI("PinnedMessages", http.MethodGet, "/channels/"+channelID.String()+"/pins",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, messages)
		})
}

// PinMessage mocks a PinMessage request.
func (m *Mocker) PinMessage(channelID, messageID discord.Snowflake) {
	m.MockAPI("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

type addRecipientPayload struct {
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nickname"`
}

// AddRecipient mocks a AddRecipient request.
func (m *Mocker) AddRecipient(channelID, userID discord.Snowflake, accessToken, nickname string) {
	m.MockAPI("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/recipients/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := addRecipientPayload{
				AccessToken: accessToken,
				Nickname:    nickname,
			}

			mockutil.CheckJSON(t, r.Body, new(addRecipientPayload), &expect)

			w.WriteHeader(http.StatusNoContent)
		})
}

// RemoveRecipient mocks a RemoveRecipient request.
func (m *Mocker) RemoveRecipient(channelID, userID discord.Snowflake) {
	m.MockAPI("RemoveRecipient", http.MethodDelete, "/channels/"+channelID.String()+"/recipients/"+userID.String(), nil)
}

// Ack mocks a Ack request.
func (m *Mocker) Ack(channelID, messageID discord.Snowflake, send, ret api.Ack) {
	m.MockAPI("Ack", http.MethodPost, "/channels/"+channelID.String()+"/messages/"+messageID.String()+"/ack",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.Ack), &send)
			mockutil.WriteJSON(t, w, ret)
		})
}
