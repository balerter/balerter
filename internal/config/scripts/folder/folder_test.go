package folder

import "testing"

func TestScriptSourceFolder_Validate(t *testing.T) {
	type fields struct {
		Name string
		Path string
		Mask string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Path: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty path",
			fields:  fields{Name: "a", Path: ""},
			wantErr: true,
			errText: "path must be not empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Folder{
				Name: tt.fields.Name,
				Path: tt.fields.Path,
				Mask: tt.fields.Mask,
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
