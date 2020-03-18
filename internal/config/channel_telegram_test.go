package config

import (
	"testing"
	"time"
)

func TestChannelTelegram_Validate(t *testing.T) {
	type fields struct {
		Name    string
		Token   string
		ChatID  int64
		Timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Token: "", ChatID: 0},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty token",
			fields:  fields{Name: "foo", Token: "", ChatID: 0},
			wantErr: true,
			errText: "token must be not empty",
		},
		{
			name:    "empty chat id",
			fields:  fields{Name: "foo", Token: "foo", ChatID: 0},
			wantErr: true,
			errText: "chat id must be not empty",
		},
		{
			name:    "low timeout",
			fields:  fields{Name: "foo", Token: "foo", ChatID: 10, Timeout: -1},
			wantErr: true,
			errText: "timeout must be greater than 0",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", Token: "foo", ChatID: 10},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelTelegram{
				Name:    tt.fields.Name,
				Token:   tt.fields.Token,
				ChatID:  tt.fields.ChatID,
				Timeout: tt.fields.Timeout,
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
