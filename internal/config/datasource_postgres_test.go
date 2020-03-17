package config

import "testing"

func TestDataSourcePostgres_Validate(t *testing.T) {
	type fields struct {
		Name        string
		Host        string
		Port        int
		Username    string
		Password    string
		Database    string
		SSLMode     string
		SSLCertPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", Host: "", Port: 0},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty host",
			fields:  fields{Name: "a", Host: "", Port: 0},
			wantErr: true,
			errText: "host must be defined",
		},
		{
			name:    "empty port",
			fields:  fields{Name: "a", Host: "a", Port: 0},
			wantErr: true,
			errText: "port must be defined",
		},
		{
			name:    "ok",
			fields:  fields{Name: "a", Host: "a", Port: 10},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DataSourcePostgres{
				Name:        tt.fields.Name,
				Host:        tt.fields.Host,
				Port:        tt.fields.Port,
				Username:    tt.fields.Username,
				Password:    tt.fields.Password,
				Database:    tt.fields.Database,
				SSLMode:     tt.fields.SSLMode,
				SSLCertPath: tt.fields.SSLCertPath,
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
