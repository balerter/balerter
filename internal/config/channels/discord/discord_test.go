package discord

import "testing"

func TestChannelDiscord_Validate(t *testing.T) {
	type fields struct {
		Name      string
		Token     string
		ChannelID int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Token: "", ChannelID: 123},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty token",
			fields:  fields{Name: "foo", Token: "", ChannelID: 123},
			wantErr: true,
			errText: "token must be not empty",
		},
		{
			name:    "empty channel",
			fields:  fields{Name: "foo", Token: "foo@bar.com", ChannelID: 0},
			wantErr: true,
			errText: "channel id must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", Token: "foo@bar.com", ChannelID: 123},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Discord{
				Name:      tt.fields.Name,
				Token:     tt.fields.Token,
				ChannelID: tt.fields.ChannelID,
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
