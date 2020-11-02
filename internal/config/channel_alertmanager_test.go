package config

import "testing"

func TestChannelAlertmanager_Validate(t *testing.T) {
	type fields struct {
		Name string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelAlertmanager{
				Name:     tt.fields.Name,
				Settings: &WebhookSettings{},
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
