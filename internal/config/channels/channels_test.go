package channels

import (
	"testing"

	"github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/balerter/balerter/internal/config/channels/log"
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/balerter/balerter/internal/config/channels/syslog"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"github.com/balerter/balerter/internal/config/channels/twiliovoice"
	"github.com/balerter/balerter/internal/config/channels/webhook"
)

func TestChannels_Validate(t *testing.T) {
	type fields struct {
		Email    []email.Email
		Slack    []slack.Slack
		Telegram []telegram.Telegram
		Syslog   []syslog.Syslog
		Notify   []notify.Notify
		Discord  []discord.Discord
		Webhook  []webhook.Webhook
		Twilio   []twiliovoice.Twilio
		Log      []log.Log
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated email",
			fields: fields{
				Email: []email.Email{{Name: "foo", From: "gopher@example.net", To: "foo@example.com",
					Host: "mail.example.com", Port: "25"}, {Name: "foo", From: "gopher@example.net",
					To: "foo@example.com", Host: "mail.example.com", Port: "25"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'email': foo",
		},
		{
			name: "duplicated slack",
			fields: fields{
				Slack: []slack.Slack{{Name: "foo", Token: "a", Channel: "a"}, {Name: "foo", Token: "a", Channel: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'slack': foo",
		},
		{
			name: "duplicated telegram",
			fields: fields{
				Telegram: []telegram.Telegram{{Name: "foo", Token: "a", ChatID: 1}, {Name: "foo", Token: "a", ChatID: 1}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'telegram': foo",
		},
		{
			name: "duplicated syslog",
			fields: fields{
				Syslog: []syslog.Syslog{{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"},
					{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'syslog': foo",
		},
		{
			name: "duplicated notify",
			fields: fields{
				Notify: []notify.Notify{{Name: "foo"}, {Name: "foo"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'notify': foo",
		},
		{
			name: "duplicated discord",
			fields: fields{
				Discord: []discord.Discord{{Name: "foo", Token: "a", ChannelID: 1}, {Name: "foo", Token: "a", ChannelID: 1}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'discord': foo",
		},
		{
			name: "duplicated webhook",
			fields: fields{
				Webhook: []webhook.Webhook{
					{
						Name: "foo",
						Settings: webhook.Settings{
							URL:     "https://foo.bar/baz",
							Auth:    &webhook.AuthConfig{Type: "bearer", AuthBearerConfig: webhook.AuthBearerConfig{Token: "token"}},
							Payload: webhook.PayloadConfig{Body: `{}`},
						},
					},
					{
						Name: "foo",
						Settings: webhook.Settings{
							URL:     "https://foo.bar/baz",
							Auth:    &webhook.AuthConfig{Type: "bearer", AuthBearerConfig: webhook.AuthBearerConfig{Token: "token"}},
							Payload: webhook.PayloadConfig{Body: `{}`},
						},
					},
				},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'webhook': foo",
		},
		{
			name: "twilio validation",
			fields: fields{
				Twilio: []twiliovoice.Twilio{{Name: "1"}},
			},
			wantErr: true,
			errText: "validate channel twilio: sid must be not empty",
		},
		{
			name: "duplicated twilio",
			fields: fields{
				Twilio: []twiliovoice.Twilio{{Name: "1", SID: "1", Token: "1", From: "1", To: "1"}, {Name: "1", SID: "1", Token: "1", From: "1", To: "1"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'twilio': 1",
		},
		{
			name: "duplicated log",
			fields: fields{
				Log: []log.Log{{Name: "1"}, {Name: "1"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'log': 1",
		},
		{
			name: "ok",
			fields: fields{
				Email: []email.Email{{Name: "foo", From: "gopher@example.net", To: "foo@example.com",
					Host: "mail.example.com", Port: "25"},
					{Name: "foo2", From: "gopher@example.net", To: "foo@example.com", Host: "mail.example.com", Port: "25"}},
				Slack: []slack.Slack{{Name: "foo", Token: "a", Channel: "a"},
					{Name: "foo2", Token: "a", Channel: "a"}},
				Telegram: []telegram.Telegram{{Name: "foo", Token: "a", ChatID: 1},
					{Name: "foo2", Token: "a", ChatID: 1}},
				Syslog: []syslog.Syslog{{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"},
					{Name: "foo2", Network: "tcp", Address: "a", Priority: "EMERG"}},
				Notify: []notify.Notify{{Name: "foo"}, {Name: "foo2"}},
				Discord: []discord.Discord{{Name: "foo", Token: "a", ChannelID: 1},
					{Name: "foo2", Token: "a", ChannelID: 1}},
				Webhook: []webhook.Webhook{
					{Name: "foo", Settings: webhook.Settings{URL: "https://foo.bar/baz",
						Auth:    &webhook.AuthConfig{Type: "bearer", AuthBearerConfig: webhook.AuthBearerConfig{Token: "token"}},
						Payload: webhook.PayloadConfig{Body: `{}`}}},
					{Name: "foo2", Settings: webhook.Settings{URL: "https://foo.bar/baz",
						Auth:    &webhook.AuthConfig{Type: "bearer", AuthBearerConfig: webhook.AuthBearerConfig{Token: "token"}},
						Payload: webhook.PayloadConfig{Body: `{}`}}}},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Channels{
				Email:       tt.fields.Email,
				Slack:       tt.fields.Slack,
				Telegram:    tt.fields.Telegram,
				Syslog:      tt.fields.Syslog,
				Notify:      tt.fields.Notify,
				Discord:     tt.fields.Discord,
				Webhook:     tt.fields.Webhook,
				TwilioVoice: tt.fields.Twilio,
				Log:         tt.fields.Log,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errText {
				t.Errorf("Validate() error = '%s', wantErrText '%s'", err.Error(), tt.errText)
			}
		})
	}
}
