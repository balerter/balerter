package postgres

import (
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"testing"
)

func TestPostgres_Validate(t *testing.T) {
	type fields struct {
		Name        string
		Host        string
		Port        int
		Timeout     int
		TableAlerts tables.TableAlerts
		TableKV     tables.TableKV
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
				Name:        "",
				Host:        "",
				Port:        0,
				Timeout:     0,
				TableAlerts: tables.TableAlerts{},
				TableKV:     tables.TableKV{},
			},
			wantErr:  true,
			errValue: "name must be not empty",
		},
		{
			name: "no host",
			fields: fields{
				Name:        "1",
				Host:        "",
				Port:        0,
				Timeout:     0,
				TableAlerts: tables.TableAlerts{},
				TableKV:     tables.TableKV{},
			},
			wantErr:  true,
			errValue: "host must be defined",
		},
		{
			name: "no port",
			fields: fields{
				Name:        "1",
				Host:        "1",
				Port:        0,
				Timeout:     0,
				TableAlerts: tables.TableAlerts{},
				TableKV:     tables.TableKV{},
			},
			wantErr:  true,
			errValue: "port must be defined",
		},
		{
			name: "bad timeout",
			fields: fields{
				Name:        "1",
				Host:        "1",
				Port:        10,
				Timeout:     -1,
				TableAlerts: tables.TableAlerts{},
				TableKV:     tables.TableKV{},
			},
			wantErr:  true,
			errValue: "timeout must be greater than 0",
		},
		{
			name: "bad TableAlerts",
			fields: fields{
				Name:    "1",
				Host:    "1",
				Port:    10,
				Timeout: 10,
				TableAlerts: tables.TableAlerts{
					Table: "",
					Fields: tables.AlertFields{
						Name:      "1",
						Level:     "1",
						Count:     "1",
						UpdatedAt: "1",
						CreatedAt: "1",
					},
					CreateTable: false,
				},
				TableKV: tables.TableKV{},
			},
			wantErr:  true,
			errValue: "table must be not empty",
		},
		{
			name: "bad TableKV",
			fields: fields{
				Name:    "1",
				Host:    "1",
				Port:    10,
				Timeout: 10,
				TableAlerts: tables.TableAlerts{
					Table: "1",
					Fields: tables.AlertFields{
						Name:      "1",
						Level:     "1",
						Count:     "1",
						UpdatedAt: "1",
						CreatedAt: "1"},
					CreateTable: false,
				},
				TableKV: tables.TableKV{},
			},
			wantErr:  true,
			errValue: "table must be not empty",
		},
		{
			name: "ok",
			fields: fields{
				Name:    "1",
				Host:    "1",
				Port:    10,
				Timeout: 10,
				TableAlerts: tables.TableAlerts{
					Table: "1",
					Fields: tables.AlertFields{
						Name:      "1",
						Level:     "1",
						Count:     "1",
						UpdatedAt: "1",
						CreatedAt: "1"},
					CreateTable: false,
				},
				TableKV: tables.TableKV{
					Table: "1",
					Fields: tables.KVFields{
						Key:   "1",
						Value: "1",
					},
					CreateTable: false,
				},
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Postgres{
				Name:        tt.fields.Name,
				Host:        tt.fields.Host,
				Port:        tt.fields.Port,
				Timeout:     tt.fields.Timeout,
				TableAlerts: tt.fields.TableAlerts,
				TableKV:     tt.fields.TableKV,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t.Errorf("unexpected error value '%s', expect '%s'", err.Error(), tt.errValue)
			}
		})
	}
}
