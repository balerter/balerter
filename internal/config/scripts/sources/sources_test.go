package sources

import (
	"github.com/balerter/balerter/internal/config/scripts/sources/folder"
	"testing"
)

func TestScriptsSources_Validate(t *testing.T) {
	type fields struct {
		Folder []folder.Folder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "empty name",
			fields: fields{
				Folder: []folder.Folder{{Name: "", Path: "", Mask: ""}},
			},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name: "empty path",
			fields: fields{
				Folder: []folder.Folder{{Name: "n1", Path: "", Mask: ""}},
			},
			wantErr: true,
			errText: "path must be not empty",
		},
		{
			name: "duplicated names",
			fields: fields{
				Folder: []folder.Folder{
					{Name: "n1", Path: ".", Mask: ""},
					{Name: "n1", Path: ".", Mask: ""},
				},
			},
			wantErr: true,
			errText: "found duplicated name for scritsource 'folder': n1",
		},
		{
			name: "ok",
			fields: fields{
				Folder: []folder.Folder{
					{Name: "n1", Path: ".", Mask: ""},
					{Name: "n2", Path: ".", Mask: ""},
				},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Sources{
				Folder: tt.fields.Folder,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("Validate() error text = %s, want err text: %s", err.Error(), tt.errText)
			}
		})
	}
}
