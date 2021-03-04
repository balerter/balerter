package webhook

import (
	"testing"
)

func TestChannelWebhook_Validate(t *testing.T) {
	type fields struct {
		Name    string
		URL     string
		Method  string
		Auth    AuthConfig
		Payload PayloadConfig
		Timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty url",
			fields:  fields{Name: "foo"},
			wantErr: true,
			errText: "url must be not empty",
		},
		{
			name:    "low timeout",
			fields:  fields{Name: "foo", URL: "https://foo.bar/baz", Timeout: -1},
			wantErr: true,
			errText: "timeout must be greater or equals 0",
		},
		{
			name: "incorrect auth type",
			fields: fields{Name: "foo", URL: "https://foo.bar/baz",
				Auth: AuthConfig{Type: "incorrect"}},
			wantErr: true,
			errText: "error validate auth: type must be set to none, basic, bearer or custom",
		},
		{
			name: "empty login",
			fields: fields{Name: "foo", URL: "https://foo.bar/baz",
				Auth: AuthConfig{Type: "basic"}},
			wantErr: true,
			errText: "error validate auth: login must be not empty",
		},
		{
			name: "empty password",
			fields: fields{Name: "foo", URL: "https://foo.bar/baz",
				Auth: AuthConfig{Type: "basic", AuthBasicConfig: AuthBasicConfig{Login: "login"}}},
			wantErr: true,
			errText: "error validate auth: password must be not empty",
		},
		{
			name: "empty token",
			fields: fields{Name: "foo", URL: "https://foo.bar/baz",
				Auth: AuthConfig{Type: "bearer"}},
			wantErr: true,
			errText: "error validate auth: token must be not empty",
		},
		{
			name:    "empty body",
			fields:  fields{Name: "foo", URL: "https://foo.bar/baz"},
			wantErr: true,
			errText: "error validate payload: body must be not empty",
		},
		{
			name: "post basic ok",
			fields: fields{
				Name: "foo",
				URL:  "https://foo.bar/baz",
				Auth: AuthConfig{
					Type:            "basic",
					AuthBasicConfig: AuthBasicConfig{Login: "login", Password: "password"},
				},
				Payload: PayloadConfig{Body: `{}`},
			},
			wantErr: false,
		},
		{
			name: "post bearer ok",
			fields: fields{
				Name: "foo",
				URL:  "https://foo.bar/baz",
				Auth: AuthConfig{
					Type:             "bearer",
					AuthBearerConfig: AuthBearerConfig{Token: "token"},
				},
				Payload: PayloadConfig{Body: `{}`},
			},
			wantErr: false,
		},
		{
			name: "post custom ok",
			fields: fields{
				Name: "foo",
				URL:  "https://foo.bar/baz",
				Auth: AuthConfig{
					Type: "custom",
					AuthCustomConfig: AuthCustomConfig{
						Headers: map[string]string{"X-Header": "value"},
					},
				},
				Payload: PayloadConfig{Body: `{}`},
			},
			wantErr: false,
		},
		{
			name: "get ok",
			fields: fields{
				Name:   "foo",
				URL:    "https://foo.bar/baz",
				Method: "get",
				Payload: PayloadConfig{
					QueryParams: map[string]string{"param": "value"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Webhook{
				Name: tt.fields.Name,
				Settings: Settings{
					URL:     tt.fields.URL,
					Method:  tt.fields.Method,
					Auth:    tt.fields.Auth,
					Payload: tt.fields.Payload,
					Timeout: tt.fields.Timeout,
				},
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
