package config

import "testing"

func TestScriptSourceFile_Validate(t *testing.T) {
	type fields struct {
		Name     string
		Filename string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Filename: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty filename",
			fields:  fields{Name: "a", Filename: ""},
			wantErr: true,
			errText: "filename must be not empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ScriptSourceFile{
				Name:     tt.fields.Name,
				Filename: tt.fields.Filename,
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
