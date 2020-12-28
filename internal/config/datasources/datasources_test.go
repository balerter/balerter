package datasources

import (
	"github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/config/datasources/loki"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"testing"
)

func TestDataSources_Validate(t *testing.T) {
	type fields struct {
		Clickhouse []*clickhouse.Clickhouse
		Prometheus []*prometheus.Prometheus
		Postgres   []*postgres.Postgres
		MySQL      []*mysql.Mysql
		Loki       []*loki.Loki
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
				Clickhouse: []*clickhouse.Clickhouse{{Name: "a", Host: "a", Port: 10}, {Name: "a", Host: "a", Port: 10}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'clickhouse': a",
		},
		{
			name: "duplicated prometheus",
			fields: fields{
				Prometheus: []*prometheus.Prometheus{{Name: "a", URL: "a"}, {Name: "a", URL: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'prometheus': a",
		},
		{
			name: "duplicated postgres",
			fields: fields{
				Postgres: []*postgres.Postgres{{Name: "a", Host: "a", Port: 10}, {Name: "a", Host: "a", Port: 10}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'postgres': a",
		},
		{
			name: "duplicated mysql",
			fields: fields{
				MySQL: []*mysql.Mysql{{Name: "a", DSN: "a"}, {Name: "a", DSN: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'mysql': a",
		},
		{
			name: "duplicated loki",
			fields: fields{
				Loki: []*loki.Loki{{Name: "a", URL: "a"}, {Name: "a", URL: "a"}},
			},
			wantErr: true,
			errText: "found duplicated name for datasource 'loki': a",
		},
		{
			name: "ok",
			fields: fields{
				Clickhouse: []*clickhouse.Clickhouse{{Name: "a", Host: "a", Port: 10}, {Name: "a2", Host: "a", Port: 10}},
				Prometheus: []*prometheus.Prometheus{{Name: "a", URL: "a"}, {Name: "a2", URL: "a"}},
				Postgres:   []*postgres.Postgres{{Name: "a", Host: "a", Port: 10}, {Name: "a2", Host: "a", Port: 10}},
				MySQL:      []*mysql.Mysql{{Name: "a", DSN: "a"}, {Name: "a2", DSN: "a"}},
				Loki:       []*loki.Loki{{Name: "a", URL: "a"}, {Name: "a2", URL: "a"}},
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
