package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/httputil/httpdriver"
	"github.com/diamondburned/arikawa/webhook"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// CreateWebhook mocks a CreateWebhook request.
//
// This method will sanitize Webhook.ID, Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) CreateWebhook(d api.CreateWebhookData, wh discord.Webhook) {
	wh = sanitize.Webhook(wh, 1, 1, 1)

	m.MockAPI("CreateWebhook", http.MethodPost, "/channels/"+wh.ChannelID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// ChannelWebhooks mocks a ChannelWebhooks request.
//
// This method will sanitize Webhook.ID, Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ChannelWebhooks(channelID discord.ChannelID, webhooks []discord.Webhook) {
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
func (m *Mocker) GuildWebhooks(guildID discord.GuildID, webhooks []discord.Webhook) {
	for i, w := range webhooks {
		webhooks[i] = sanitize.Webhook(w, 1, 1, 1)

		if w.GuildID == 0 {
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
// The ID field of the passed discord.Webhook must be set.
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
// The ID field and the Token field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) WebhookWithToken(wh discord.Webhook) {
	wh = sanitize.Webhook(wh, 1, 1, 1)

	webhook.DefaultHTTPClient.Client = httpdriver.WrapClient(*m.Client)

	m.MockAPI("WebhookWithToken", http.MethodGet, "/webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, wh)
		})
}

// ModifyWebhook mocks a ModifyWebhook request.
//
// The ID field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhook(d api.ModifyWebhookData, wh discord.Webhook) {
	wh = sanitize.Webhook(wh, 1, 1, 1)

	m.MockAPI("ModifyWebhook", http.MethodPatch, "/webhooks/"+wh.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// ModifyWebhookWithToken mocks a ModifyWebhookWithToken request.
//
// The ID field and the Token field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhookWithToken(d api.ModifyWebhookData, wh discord.Webhook) {
	wh = sanitize.Webhook(wh, 1, 1, 1)

	webhook.DefaultHTTPClient.Client = httpdriver.WrapClient(*m.Client)

	m.MockAPI("ModifyWebhookWithToken", http.MethodPatch, "/webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// DeleteWebhook mocks a DeleteWebhook request.
func (m *Mocker) DeleteWebhook(id discord.WebhookID) {
	m.MockAPI("DeleteWebhook", http.MethodDelete, "/webhooks/"+id.String(), nil)
}

// DeleteWebhookWithToken mocks a DeleteWebhookWithToken request.
func (m *Mocker) DeleteWebhookWithToken(id discord.WebhookID, token string) {
	webhook.DefaultHTTPClient.Client = httpdriver.WrapClient(*m.Client)

	m.MockAPI("DeleteWebhookWithToken", http.MethodDelete, "/webhooks/"+id.String()+"/"+token, nil)
}
