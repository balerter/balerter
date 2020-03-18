package config

import (
	"testing"
	"time"
)

func TestDataSourceMysql_Validate(t *testing.T) {
	type fields struct {
		Name    string
		DSN     string
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
			fields:  fields{Name: "", DSN: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty dsn",
			fields:  fields{Name: "foo", DSN: ""},
			wantErr: true,
			errText: "DSN must be not empty",
		},
		{
			name:    "wrong timeout",
			fields:  fields{Name: "foo", DSN: "a", Timeout: -1},
			wantErr: true,
			errText: "timeout must be greater than 0",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", DSN: "a"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DataSourceMysql{
				Name:    tt.fields.Name,
				DSN:     tt.fields.DSN,
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
