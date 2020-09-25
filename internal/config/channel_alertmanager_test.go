package config

import "testing"

func TestChannelAlertmanager_Validate(t *testing.T) {
	type fields struct {
		Name    string
		Version string
		URL     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Version: "", URL: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty url",
			fields:  fields{Name: "foo", Version: "", URL: ""},
			wantErr: true,
			errText: "url must be not empty",
		},
		{
			name:    "wrong version",
			fields:  fields{Name: "foo", Version: "bar", URL: "foo"},
			wantErr: true,
			errText: "version must be empty or v1 or v2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelAlertmanager{
				Name:    tt.fields.Name,
				Version: tt.fields.Version,
				URL:     tt.fields.URL,
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
