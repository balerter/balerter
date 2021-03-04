package sqlite

import (
	"testing"
)

func TestStorageCoreFile_Validate(t *testing.T) {
	type fields struct {
		Name        string
		Path        string
		Timeout     int
		TableAlerts string
		TableKV     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Path: "", Timeout: -1},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty path",
			fields:  fields{Name: "a", Path: "", Timeout: -1},
			wantErr: true,
			errText: "path must be not empty",
		},
		{
			name:    "empty table alerts",
			fields:  fields{Name: "a", Path: "a", Timeout: -1, TableAlerts: "", TableKV: ""},
			wantErr: true,
			errText: "table Alerts must be not empty",
		},
		{
			name:    "empty table kv",
			fields:  fields{Name: "a", Path: "a", Timeout: -1, TableAlerts: "foo", TableKV: ""},
			wantErr: true,
			errText: "table KV must be not empty",
		},
		{
			name:    "empty timeout",
			fields:  fields{Name: "a", Path: "a", Timeout: -1, TableAlerts: "foo", TableKV: "bar"},
			wantErr: true,
			errText: "timeout must be greater than 0",
		},
		{
			name:    "ok",
			fields:  fields{Name: "a", Path: "a", Timeout: 10, TableAlerts: "foo", TableKV: "bar"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Sqlite{
				Name:    tt.fields.Name,
				Path:    tt.fields.Path,
				Timeout: tt.fields.Timeout,
				Tables: Tables{
					Alerts: tt.fields.TableAlerts,
					KV:     tt.fields.TableKV,
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
