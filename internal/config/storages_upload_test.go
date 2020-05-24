package config

import "testing"

func TestStoragesUpload_Validate(t *testing.T) {
	type fields struct {
		S3 []*StorageUploadS3
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated s3",
			fields: fields{
				S3: []*StorageUploadS3{{Name: "a"}, {Name: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for upload storages 's3': a",
		},
		{
			name: "ok",
			fields: fields{
				S3: []*StorageUploadS3{{Name: "a"}, {Name: "a2"}},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := StoragesUpload{
				S3: tt.fields.S3,
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
