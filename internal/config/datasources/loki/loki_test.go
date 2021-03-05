package loki

import (
	"github.com/balerter/balerter/internal/config/common"
	"testing"
)

func TestDataSourceLoki_Validate(t *testing.T) {
	type fields struct {
		Name      string
		URL       string
		BasicAuth *common.BasicAuth
		Timeout   int
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
			name:    "empty URL",
			fields:  fields{Name: "foo", URL: ""},
			wantErr: true,
			errText: "url must be not empty",
		},
		{
			name:    "wrong timeout",
			fields:  fields{Name: "foo", URL: "a", Timeout: -1},
			wantErr: true,
			errText: "timeout must be greater than 0",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", URL: "a"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Loki{
				Name:      tt.fields.Name,
				URL:       tt.fields.URL,
				BasicAuth: tt.fields.BasicAuth,
				Timeout:   tt.fields.Timeout,
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
