package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// Emojis mocks a Emojis request.
//
// This method will sanitize Emoji.ID and Emoji.User.ID.
func (m *Mocker) Emojis(guildID discord.Snowflake, e []discord.Emoji) {
	if e == nil {
		e = []discord.Emoji{}
	}

	for i, emoji := range e {
		e[i] = sanitize.Emoji(emoji, 1, 1)
	}

	m.MockAPI("Emojis", http.MethodGet, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// Emoji mocks a Emoji request.
//
// The ID field of the passed Emoji is required.
//
// This method will sanitize Emoji.ID and Emoji.User.ID.
func (m *Mocker) Emoji(guildID discord.Snowflake, e discord.Emoji) {
	e = sanitize.Emoji(e, 1, 1)

	m.MockAPI("Emoji", http.MethodGet, "/guilds/"+guildID.String()+"/emojis/"+e.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// CreateEmoji mocks a CreateEmoji request.
//
// The fields Name and RoleIDs of the passed Emoji must be set.
//
// This method will sanitize Emoji.ID and Emoji.User.ID.
func (m *Mocker) CreateEmoji(guildID discord.Snowflake, d api.CreateEmojiData, e discord.Emoji) {
	e = sanitize.Emoji(e, 1, 1)

	m.MockAPI("CreateEmoji", http.MethodPost, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateEmojiData), &d)
			mockutil.WriteJSON(t, w, e)
		})
}

// ModifyEmoji mocks a ModifyEmoji request.
func (m *Mocker) ModifyEmoji(guildID, emojiID discord.Snowflake, d api.ModifyEmojiData) {
	m.MockAPI("ModifyEmoji", http.MethodPatch, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyEmojiData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteEmoji mocks a DeleteEmoji request.
func (m *Mocker) DeleteEmoji(guildID, emojiID discord.Snowflake) {
	m.MockAPI("DeleteEmoji", http.MethodDelete, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(), nil)
}
