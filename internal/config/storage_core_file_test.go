package config

import (
	"testing"
	"time"
)

func TestStorageCoreFile_Validate(t *testing.T) {
	type fields struct {
		Name    string
		Path    string
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
			name:    "empty timeout",
			fields:  fields{Name: "a", Path: "a", Timeout: -1},
			wantErr: true,
			errText: "timeout must be greater than 0",
		},
		{
			name:    "ok",
			fields:  fields{Name: "a", Path: "a", Timeout: 10},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := StorageCoreFile{
				Name:    tt.fields.Name,
				Path:    tt.fields.Path,
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
