package config

import "testing"

func TestDataSources_Validate(t *testing.T) {
	type fields struct {
		Clickhouse []DataSourceClickhouse
		Prometheus []DataSourcePrometheus
		Postgres   []DataSourcePostgres
		MySQL      []DataSourceMysql
		Loki       []DataSourceLoki
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated clickhouse",
			fields: fields{
				Clickhouse: []DataSourceClickhouse{{Name: "a", Host: "a", Port: 10}, {Name: "a", Host: "a", Port: 10}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'clickhouse': a",
		},
		{
			name: "duplicated prometheus",
			fields: fields{
				Prometheus: []DataSourcePrometheus{{Name: "a", URL: "a"}, {Name: "a", URL: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'prometheus': a",
		},
		{
			name: "duplicated postgres",
			fields: fields{
				Postgres: []DataSourcePostgres{{Name: "a", Host: "a", Port: 10}, {Name: "a", Host: "a", Port: 10}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'postgres': a",
		},
		{
			name: "duplicated mysql",
			fields: fields{
				MySQL: []DataSourceMysql{{Name: "a", DSN: "a"}, {Name: "a", DSN: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'mysql': a",
		},
		{
			name: "duplicated loki",
			fields: fields{
				Loki: []DataSourceLoki{{Name: "a", URL: "a"}, {Name: "a", URL: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'loki': a",
		},
		{
			name: "ok",
			fields: fields{
				Clickhouse: []DataSourceClickhouse{{Name: "a", Host: "a", Port: 10}, {Name: "a2", Host: "a", Port: 10}},
				Prometheus: []DataSourcePrometheus{{Name: "a", URL: "a"}, {Name: "a2", URL: "a"}},
				Postgres:   []DataSourcePostgres{{Name: "a", Host: "a", Port: 10}, {Name: "a2", Host: "a", Port: 10}},
				MySQL:      []DataSourceMysql{{Name: "a", DSN: "a"}, {Name: "a2", DSN: "a"}},
				Loki:       []DataSourceLoki{{Name: "a", URL: "a"}, {Name: "a2", URL: "a"}},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DataSources{
				Clickhouse: tt.fields.Clickhouse,
				Prometheus: tt.fields.Prometheus,
				Postgres:   tt.fields.Postgres,
				MySQL:      tt.fields.MySQL,
				Loki:       tt.fields.Loki,
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
