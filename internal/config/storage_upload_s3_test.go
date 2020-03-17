package config

import "testing"

func TestStorageUploadS3_Validate(t *testing.T) {
	type fields struct {
		Name     string
		Region   string
		Key      string
		Secret   string
		Endpoint string
		Bucket   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "a"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := StorageUploadS3{
				Name:     tt.fields.Name,
				Region:   tt.fields.Region,
				Key:      tt.fields.Key,
				Secret:   tt.fields.Secret,
				Endpoint: tt.fields.Endpoint,
				Bucket:   tt.fields.Bucket,
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
