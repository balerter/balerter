package slack

import "testing"

func TestChannelSlack_Validate(t *testing.T) {
	type fields struct {
		Name    string
		Token   string
		Channel string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Token: "", Channel: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty token",
			fields:  fields{Name: "foo", Token: "", Channel: ""},
			wantErr: true,
			errText: "token must be not empty",
		},
		{
			name:    "empty channel",
			fields:  fields{Name: "foo", Token: "foo", Channel: ""},
			wantErr: true,
			errText: "channel must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", Token: "foo", Channel: "foo"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Slack{
				Name:    tt.fields.Name,
				Token:   tt.fields.Token,
				Channel: tt.fields.Channel,
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
