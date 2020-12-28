package core

import (
	"github.com/balerter/balerter/internal/config/storages/core/file"
	"testing"
)

func TestStoragesCore_Validate(t *testing.T) {
	type fields struct {
		File []*file.File
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated file",
			fields: fields{
				File: []*file.File{{Name: "a", Path: "a"}, {Name: "a", Path: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for core storages 'file': a",
		},
		{
			name: "ok",
			fields: fields{
				File: []*file.File{{Name: "a", Path: "a"}, {Name: "a2", Path: "a"}},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Core{
				File: tt.fields.File,
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
