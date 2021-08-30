package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/common"
	"github.com/balerter/balerter/internal/config/datasources"
	clickhouseCfg "github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/config/datasources/loki"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestManager_Stop(t *testing.T) {
	m1 := &modules.ModuleMock{
		StopFunc: func() error {
			return nil
		},
		NameFunc: func() string {
			return "m1"
		},
	}

	m2 := &modules.ModuleMock{
		StopFunc: func() error {
			return nil
		},
		NameFunc: func() string {
			return "m2"
		},
	}

	m := New(zap.NewNop())
	m.modules["m1"] = m1
	m.modules["m2"] = m2

	m.Stop()

	assert.Equal(t, 1, len(m1.StopCalls()))
	assert.Equal(t, 1, len(m1.NameCalls()))
	assert.Equal(t, 1, len(m2.StopCalls()))
	assert.Equal(t, 1, len(m2.NameCalls()))
}

func TestManager_Stop_Error(t *testing.T) {
	e := fmt.Errorf("error1")

	m1 := &modules.ModuleMock{
		StopFunc: func() error {
			return e
		},
		NameFunc: func() string {
			return "m1"
		},
	}

	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	m := New(logger)
	m.modules["m1"] = m1

	m.Stop()

	assert.Equal(t, 1, logs.FilterMessage("error stop module").FilterField(zap.String("name", "m1")).FilterField(zap.Error(e)).Len())

	assert.Equal(t, 1, len(m1.StopCalls()))
	assert.Equal(t, 2, len(m1.NameCalls()))
}

func TestGet(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.Module{
			"m1": &modules.ModuleMock{},
			"m2": &modules.ModuleMock{},
		},
	}

	mods := m.Get()

	assert.Equal(t, 2, len(mods))
}

func TestManager_Init(t *testing.T) {
	type fields struct {
		modules map[string]modules.Module
	}
	type args struct {
		cfg *datasources.DataSources
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil config",
			fields: fields{
				modules: nil,
			},
			args:    args{cfg: nil},
			wantErr: false,
		},
		{
			name: "err clickhouse",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					Clickhouse: []clickhouseCfg.Clickhouse{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "err prometheus",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					Prometheus: []prometheus.Prometheus{{URL: "://://"}},
				},
			},
			wantErr: true,
		},
		{
			name: "err postgres",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					Postgres: []postgres.Postgres{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "err mysql",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					MySQL: []mysql.Mysql{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "err loki",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					Loki: []loki.Loki{{URL: "://://"}},
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				modules: map[string]modules.Module{},
			},
			args: args{
				cfg: &datasources.DataSources{
					Clickhouse: []clickhouseCfg.Clickhouse{{
						Name:        "foo",
						Host:        "127.0.0.1",
						Port:        9000,
						Username:    "default",
						Password:    "",
						Database:    "default",
						SSLCertPath: "",
						Timeout:     0,
					}},
					Prometheus: []prometheus.Prometheus{{
						Name:      "foo",
						URL:       "http://127.0.0.1",
						BasicAuth: &common.BasicAuth{},
						Timeout:   0,
					}},
					Postgres: []postgres.Postgres{{
						Name:        "foo",
						Host:        "127.0.0.1",
						Port:        35432,
						Username:    "postgres",
						Password:    "secret",
						Database:    "db",
						SSLMode:     "disable",
						SSLCertPath: "",
						Timeout:     0,
					}},
					MySQL: []mysql.Mysql{{
						Name:    "foo",
						DSN:     "mysql:secret@tcp(127.0.0.1:3306)/db",
						Timeout: 0,
					}},
					Loki: []loki.Loki{{
						Name:      "foo",
						URL:       "http://127.0.0.1",
						BasicAuth: nil,
						Timeout:   0,
					}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				modules: tt.fields.modules,
			}
			err := m.Init(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			//if err != nil && tt.wantErr && err.Error() != tt.errValue {
			//	t.Errorf("unexpected error value = %s, want %s", err.Error(), tt.errValue)
			//}
		})
	}
}
