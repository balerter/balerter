package config

import "testing"

func TestChannels_Validate(t *testing.T) {
	type fields struct {
		Email    []*ChannelEmail
		Slack    []*ChannelSlack
		Telegram []*ChannelTelegram
		Syslog   []*ChannelSyslog
		Notify   []*ChannelNotify
		Discord  []*ChannelDiscord
		Webhook  []*ChannelWebhook
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
				Email: []*ChannelEmail{{Name: "foo", From: "gopher@example.net", To: "foo@example.com",
					Host: "mail.example.com", Port: "25"}, {Name: "foo", From: "gopher@example.net",
					To: "foo@example.com", Host: "mail.example.com", Port: "25"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'email': foo",
		},
		{
			name: "duplicated slack",
			fields: fields{
				Slack: []*ChannelSlack{{Name: "foo", Token: "a", Channel: "a"}, {Name: "foo", Token: "a", Channel: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'slack': foo",
		},
		{
			name: "duplicated telegram",
			fields: fields{
				Telegram: []*ChannelTelegram{{Name: "foo", Token: "a", ChatID: 1}, {Name: "foo", Token: "a", ChatID: 1}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'telegram': foo",
		},
		{
			name: "duplicated syslog",
			fields: fields{
				Syslog: []*ChannelSyslog{{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"},
					{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'syslog': foo",
		},
		{
			name: "duplicated notify",
			fields: fields{
				Notify: []*ChannelNotify{{Name: "foo"}, {Name: "foo"}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'notify': foo",
		},
		{
			name: "duplicated discord",
			fields: fields{
				Discord: []*ChannelDiscord{{Name: "foo", Token: "a", ChannelID: 1}, {Name: "foo", Token: "a", ChannelID: 1}},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'discord': foo",
		},
		{
			name: "duplicated webhook",
			fields: fields{
				Webhook: []*ChannelWebhook{
					{
						Name:    "foo",
						URL:     "https://foo.bar/baz",
						Auth:    AuthConfig{Type: "bearer", AuthBearerConfig: AuthBearerConfig{Token: "token"}},
						Payload: PayloadConfig{Body: `{}`},
					},
					{
						Name:    "foo",
						URL:     "https://foo.bar/baz",
						Auth:    AuthConfig{Type: "bearer", AuthBearerConfig: AuthBearerConfig{Token: "token"}},
						Payload: PayloadConfig{Body: `{}`},
					},
				},
			},
			wantErr: true,
			errText: "found duplicated name for channels 'webhook': foo",
		},
		{
			name: "ok",
			fields: fields{
				Email: []*ChannelEmail{{Name: "foo", From: "gopher@example.net", To: "foo@example.com",
					Host: "mail.example.com", Port: "25"},
					{Name: "foo2", From: "gopher@example.net", To: "foo@example.com", Host: "mail.example.com", Port: "25"}},
				Slack: []*ChannelSlack{{Name: "foo", Token: "a", Channel: "a"},
					{Name: "foo2", Token: "a", Channel: "a"}},
				Telegram: []*ChannelTelegram{{Name: "foo", Token: "a", ChatID: 1},
					{Name: "foo2", Token: "a", ChatID: 1}},
				Syslog: []*ChannelSyslog{{Name: "foo", Network: "tcp", Address: "a", Priority: "EMERG"},
					{Name: "foo2", Network: "tcp", Address: "a", Priority: "EMERG"}},
				Notify: []*ChannelNotify{{Name: "foo"}, {Name: "foo2"}},
				Discord: []*ChannelDiscord{{Name: "foo", Token: "a", ChannelID: 1},
					{Name: "foo2", Token: "a", ChannelID: 1}},
				Webhook: []*ChannelWebhook{
					{Name: "foo", URL: "https://foo.bar/baz",
						Auth:    AuthConfig{Type: "bearer", AuthBearerConfig: AuthBearerConfig{Token: "token"}},
						Payload: PayloadConfig{Body: `{}`}},
					{Name: "foo2", URL: "https://foo.bar/baz",
						Auth:    AuthConfig{Type: "bearer", AuthBearerConfig: AuthBearerConfig{Token: "token"}},
						Payload: PayloadConfig{Body: `{}`}}},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Channels{
				Email:    tt.fields.Email,
				Slack:    tt.fields.Slack,
				Telegram: tt.fields.Telegram,
				Syslog:   tt.fields.Syslog,
				Notify:   tt.fields.Notify,
				Discord:  tt.fields.Discord,
				Webhook:  tt.fields.Webhook,
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
