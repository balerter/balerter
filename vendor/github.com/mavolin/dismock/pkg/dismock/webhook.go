package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// CreateWebhook mocks a CreateWebhook request.
//
// This method will sanitize Webhook.ID, Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) CreateWebhook(d api.CreateWebhookData, webhook discord.Webhook) {
	webhook = sanitize.Webhook(webhook, 1, 1, 1)

	m.MockAPI("CreateWebhook", http.MethodPost, "/channels/"+webhook.ChannelID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateWebhookData), &d)
			mockutil.WriteJSON(t, w, webhook)
		})
}

// ChannelWebhooks mocks a ChannelWebhooks request.
//
// This method will sanitize Webhook.ID, Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ChannelWebhooks(channelID discord.Snowflake, webhooks []discord.Webhook) {
	for i, w := range webhooks {
		webhooks[i] = sanitize.Webhook(w, 1, 1, channelID)
	}

	m.MockAPI("ChannelWebhooks", http.MethodGet, "/channels/"+channelID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhooks)
		})
}

// GuildWebhooks mocks a GuildWebhooks request.
//
// This method will sanitize Webhook.ID, Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) GuildWebhooks(guildID discord.Snowflake, webhooks []discord.Webhook) {
	for i, w := range webhooks {
		webhooks[i] = sanitize.Webhook(w, 1, 1, 1)

		if w.GuildID <= 0 {
			webhooks[i].GuildID = guildID
		}
	}

	m.MockAPI("GuildWebhooks", http.MethodGet, "/guilds/"+guildID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhooks)
		})
}

// Webhook mocks a Webhook request.
//
// The ID field of the passed webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) Webhook(webhook discord.Webhook) {
	webhook = sanitize.Webhook(webhook, 1, 1, 1)

	m.MockAPI("Webhook", http.MethodGet, "/webhooks/"+webhook.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhook)
		})
}

// WebhookWithToken mocks a WebhookWithToken request.
//
// The ID field and the Token field of the passed webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) WebhookWithToken(webhook discord.Webhook) {
	webhook = sanitize.Webhook(webhook, 1, 1, 1)

	m.MockAPI("WebhookWithToken", http.MethodGet, "/webhooks/"+webhook.ID.String()+"/"+webhook.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhook)
		})
}

// ModifyWebhook mocks a ModifyWebhook request.
//
// The ID field of the passed Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhook(d api.ModifyWebhookData, webhook discord.Webhook) {
	webhook = sanitize.Webhook(webhook, 1, 1, 1)

	m.MockAPI("ModifyWebhook", http.MethodPatch, "/webhooks/"+webhook.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, webhook)
		})
}

// ModifyWebhookWithToken mocks a ModifyWebhookWithToken request.
//
// The ID field and the Token field of the passed Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhookWithToken(d api.ModifyWebhookData, webhook discord.Webhook) {
	webhook = sanitize.Webhook(webhook, 1, 1, 1)

	m.MockAPI("ModifyWebhookWithToken", http.MethodPatch, "/webhooks/"+webhook.ID.String()+"/"+webhook.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, webhook)
		})
}

// DeleteWebhook mocks a DeleteWebhook request.
func (m *Mocker) DeleteWebhook(id discord.Snowflake) {
	m.MockAPI("DeleteWebhook", http.MethodDelete, "/webhooks/"+id.String(), nil)
}

// DeleteWebhookWithToken mocks a DeleteWebhookWithToken request.
func (m *Mocker) DeleteWebhookWithToken(id discord.Snowflake, token string) {
	m.MockAPI("DeleteWebhookWithToken", http.MethodDelete, "/webhooks/"+id.String()+"/"+token, nil)
}
