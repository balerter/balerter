package config

import "testing"

func TestChannelEmail_Validate(t *testing.T) {
	type fields struct {
		Name string
		From string
		To   string
		Host string
		Port string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", From: "", To: "", Host: "", Port: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty from",
			fields:  fields{Name: "foo", From: "", To: "", Host: "", Port: ""},
			wantErr: true,
			errText: "from must be not empty",
		},
		{
			name:    "empty to",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "", Host: "", Port: ""},
			wantErr: true,
			errText: "to must be not empty",
		},
		{
			name:    "empty host",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "foo@example.com", Host: "", Port: ""},
			wantErr: true,
			errText: "host must be not empty",
		},
		{
			name:    "empty port",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "foo@example.com", Host: "mail.example.com", Port: ""},
			wantErr: true,
			errText: "port must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "foo@example.com", Host: "mail.example.com", Port: "25"},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelEmail{
				Name: tt.fields.Name,
				From: tt.fields.From,
				To:   tt.fields.To,
				Host: tt.fields.Host,
				Port: tt.fields.Port,
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
