package dismock

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// channel is a discord.Channel from a senders perspective.
type channel struct {
	discord.Channel
	Permissions []overwrite `json:"permission_overwrites,omitempty"`
}

func newChannel(c discord.Channel) (wrapped channel) {
	wrapped = channel{
		Channel: c,
	}

	for _, p := range c.Permissions {
		o := newOverwrite(p)
		wrapped.Permissions = append(wrapped.Permissions, o)
	}

	return
}

// overwrite is the discord.Overwrite struct from a senders perspective.
type overwrite struct {
	ID    discord.Snowflake     `json:"id"`
	Type  discord.OverwriteType `json:"type"`
	Allow discord.Permissions   `json:"allow_new,string"`
	Deny  discord.Permissions   `json:"deny_new,string"`
}

func newOverwrite(o discord.Overwrite) overwrite {
	return overwrite{
		ID:    o.ID,
		Type:  o.Type,
		Allow: o.Allow,
		Deny:  o.Deny,
	}
}

func (o *overwrite) UnmarshalJSON(data []byte) (err error) {
	var recv struct {
		ID    discord.Snowflake     `json:"id"`
		Type  discord.OverwriteType `json:"type"`
		Allow discord.Permissions   `json:"allow,string"`
		Deny  discord.Permissions   `json:"deny,string"`
	}

	err = json.Unmarshal(data, &recv)
	if err != nil {
		return
	}

	o.ID = recv.ID
	o.Type = recv.Type
	o.Allow = recv.Allow
	o.Deny = recv.Deny

	return
}

// Channels mocks a channels request.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channels(guildID discord.GuildID, c []discord.Channel) {
	wrapped := make([]channel, len(c))

	for i, c := range c {
		c = sanitize.Channel(c, 1)
		wrapped[i] = newChannel(c)
	}

	m.MockAPI("Channels", http.MethodGet, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, wrapped)
		})
}

// createChannelData is a api.CreateChannelData struct from a senders
// perspective.
type createChannelData struct {
	api.CreateChannelData
	Permissions []overwrite `json:"permission_overwrites,omitempty"`
}

func newCreateChannelData(d api.CreateChannelData) (wrapped createChannelData) {
	wrapped = createChannelData{
		CreateChannelData: d,
	}

	for _, p := range d.Permissions {
		o := newOverwrite(p)
		wrapped.Permissions = append(wrapped.Permissions, o)
	}

	wrapped.CreateChannelData.Permissions = nil

	return
}

// CreateChannel mocks a CreateChannel request.
//
// The GuildID field of the passed discord.Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) CreateChannel(d api.CreateChannelData, c discord.Channel) {
	c = sanitize.Channel(c, 1)

	recv := newCreateChannelData(d)
	send := newChannel(c)

	m.MockAPI("CreateChannel", http.MethodPost, "/guilds/"+c.GuildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(createChannelData), &recv)
			mockutil.WriteJSON(t, w, send)
		})
}

// MoveChannel mocks a MoveChannel request.
func (m *Mocker) MoveChannel(guildID discord.GuildID, d []api.MoveChannelData) {
	m.MockAPI("CreateChannel", http.MethodPatch, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &[]api.MoveChannelData{}, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Channel mocks a Channel request.
//
// The ID field of the passed discord.Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channel(c discord.Channel) {
	c = sanitize.Channel(c, 1)

	send := newChannel(c)

	m.MockAPI("Channel", http.MethodGet, "/channels/"+c.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, send)
		})
}

type modifyChannelData struct {
	api.ModifyChannelData
	Permissions *[]overwrite `json:"permission_overwrites,omitempty"`
}

func newModifyChannelData(d api.ModifyChannelData) (wrapped modifyChannelData) {
	wrapped = modifyChannelData{
		ModifyChannelData: d,
	}

	if d.Permissions != nil {
		perms := make([]overwrite, len(*d.Permissions))

		for i, p := range *d.Permissions {
			o := newOverwrite(p)
			perms[i] = o
		}

		wrapped.Permissions = &perms
	}

	wrapped.ModifyChannelData.Permissions = nil

	return
}

// ModifyChannel mocks a ModifyChannel request.
func (m *Mocker) ModifyChannel(id discord.ChannelID, d api.ModifyChannelData) {
	recv := newModifyChannelData(d)

	m.MockAPI("ModifyChannel", http.MethodPatch, "/channels/"+id.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(modifyChannelData), &recv)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannel mocks a DeleteChannel request.
func (m *Mocker) DeleteChannel(id discord.ChannelID) {
	m.MockAPI("DeleteChannel", http.MethodDelete, "/channels/"+id.String(), nil)
}

// EditChannelPermission mocks a EditChannelPermission request.
func (m *Mocker) EditChannelPermission(
	channelID discord.ChannelID, overwriteID discord.Snowflake, d api.EditChannelPermissionData,
) {
	m.MockAPI("EditChannelPermission", http.MethodPut,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.EditChannelPermissionData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannelPermission mocks a DeleteChannelPermission request.
func (m *Mocker) DeleteChannelPermission(channelID discord.ChannelID, overwriteID discord.Snowflake) {
	m.MockAPI("DeleteChannelPermission", http.MethodDelete,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(), nil)
}

// Typing mocks a Typing request.
func (m *Mocker) Typing(channelID discord.ChannelID) {
	m.MockAPI("Typing", http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
}

// PinnedMessages mocks a PinnedMessages request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) PinnedMessages(channelID discord.ChannelID, messages []discord.Message) {
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
func (m *Mocker) PinMessage(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

// UnpinMessage mocks a UnpinMessage request.
func (m *Mocker) UnpinMessage(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("UnpinMessage", http.MethodDelete, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

type addRecipientPayload struct {
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nickname"`
}

// AddRecipient mocks a AddRecipient request.
func (m *Mocker) AddRecipient(channelID discord.ChannelID, userID discord.UserID, accessToken, nickname string) {
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
func (m *Mocker) RemoveRecipient(channelID discord.ChannelID, userID discord.UserID) {
	m.MockAPI("RemoveRecipient", http.MethodDelete, "/channels/"+channelID.String()+"/recipients/"+userID.String(), nil)
}

// Ack mocks a Ack request.
func (m *Mocker) Ack(channelID discord.ChannelID, messageID discord.MessageID, send, ret api.Ack) {
	m.MockAPI("Ack", http.MethodPost, "/channels/"+channelID.String()+"/messages/"+messageID.String()+"/ack",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.Ack), &send)
			mockutil.WriteJSON(t, w, ret)
		})
}
