package core

import (
	"github.com/balerter/balerter/internal/config/storages/core/sqlite"
	"testing"
)

func TestStoragesCore_Validate(t *testing.T) {
	type fields struct {
		File []sqlite.Sqlite
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated storage name",
			fields: fields{
				File: []sqlite.Sqlite{
					{Name: "a", Path: "a", Timeout: 0, Tables: sqlite.Tables{Alerts: "a", KV: "k"}},
					{Name: "a", Path: "a", Timeout: 0, Tables: sqlite.Tables{Alerts: "a", KV: "k"}},
				},
			},
			wantErr: true,
			errText: "found duplicated name for core storages 'sqlite': a",
		},
		{
			name: "ok",
			fields: fields{
				File: []sqlite.Sqlite{
					{Name: "a", Path: "a", Timeout: 0, Tables: sqlite.Tables{Alerts: "a", KV: "k"}},
					{Name: "a2", Path: "a", Timeout: 0, Tables: sqlite.Tables{Alerts: "a", KV: "k"}},
				},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Core{
				Sqlite: tt.fields.File,
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
