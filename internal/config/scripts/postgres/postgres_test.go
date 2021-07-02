package postgres

import "testing"

func TestPostgres_Validate(t *testing.T) {
	type fields struct {
		Name    string
		Host    string
		Port    int
		Timeout int
		Query   string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name: "no name",
			fields: fields{
				Name:    "",
				Host:    "",
				Port:    0,
				Timeout: 0,
				Query:   "",
			},
			wantErr:  true,
			errValue: "name must be not empty",
		},
		{
			name: "no host",
			fields: fields{
				Name:    "foo",
				Host:    "",
				Port:    0,
				Timeout: 0,
				Query:   "",
			},
			wantErr:  true,
			errValue: "host must be defined",
		},
		{
			name: "no port",
			fields: fields{
				Name:    "foo",
				Host:    "foo",
				Port:    0,
				Timeout: 0,
				Query:   "",
			},
			wantErr:  true,
			errValue: "port must be defined",
		},
		{
			name: "bad timeout",
			fields: fields{
				Name:    "foo",
				Host:    "foo",
				Port:    1000,
				Timeout: -1,
				Query:   "",
			},
			wantErr:  true,
			errValue: "timeout must be greater than 0",
		},
		{
			name: "empty query",
			fields: fields{
				Name:    "foo",
				Host:    "foo",
				Port:    1000,
				Timeout: 1000,
				Query:   "",
			},
			wantErr:  true,
			errValue: "query must be defined",
		},
		{
			name: "ok",
			fields: fields{
				Name:    "foo",
				Host:    "foo",
				Port:    1000,
				Timeout: 1000,
				Query:   "foo",
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Postgres{
				Name:    tt.fields.Name,
				Host:    tt.fields.Host,
				Port:    tt.fields.Port,
				Timeout: tt.fields.Timeout,
				Query:   tt.fields.Query,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.wantErr && tt.errValue != err.Error() {
				t.Errorf("unexpected error value %s, expect %s", err.Error(), tt.errValue)
			}
		})
	}
}
