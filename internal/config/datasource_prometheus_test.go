package config

import "testing"

func TestDataSourcePrometheus_Validate(t *testing.T) {
	type fields struct {
		Name      string
		URL       string
		BasicAuth BasicAuth
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", URL: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty url",
			fields:  fields{Name: "a", URL: ""},
			wantErr: true,
			errText: "url must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "a", URL: "a"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DataSourcePrometheus{
				Name:      tt.fields.Name,
				URL:       tt.fields.URL,
				BasicAuth: tt.fields.BasicAuth,
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
