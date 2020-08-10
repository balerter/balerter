package main

import (
	"bytes"
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func Test_Main(t *testing.T) {
//	version = "TEST"
//	os.Args = []string{"-debug", "-once", "-config", "stdin", "-script", "test.lua"}
//
//	buf := &bytes.Buffer{}
//	ws := zapcore.AddSync(buf)
//
//	loggerOptions = append(loggerOptions, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
//		return zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), ws, zap.DebugLevel)
//	}))
//
//	main()
//
//	assert.Contains(t, buf.String(), `,"msg":"balerter start","version":"TEST"}`)
//	assert.Contains(t, buf.String(), `,"msg":"Foo Bar","scriptName":"test"}`)
//}

func Test_Main_bad_log_level(t *testing.T) {
	msg, code := run("config.yml", "FOO", false, false, "")

	assert.Equal(t, "wrong log level", msg)
	assert.Equal(t, 1, code)
}

func Test_Main_error_load_config(t *testing.T) {
	msg, code := run("notfound.yml", "INFO", false, false, "")

	assert.Equal(t, "error init config, error read config file, open notfound.yml: no such file or directory", msg)
	assert.Equal(t, 1, code)
}

func Test_Main_error_init_datasource_manager(t *testing.T) {
	buf := []byte(`
datasources:
 clickhouse:
   - name: ch1
     host: 127.0.0.1
`)

	config.StdIn = bytes.NewBuffer(buf)

	msg, code := run("stdin", "INFO", false, false, "")

	assert.Equal(t, "error init config, error config validation, error DataSources validation, port must be defined", msg)
	assert.Equal(t, 1, code)
}

func Test_Main_error_init_storages_alert(t *testing.T) {
	buf := []byte(`
global:
 storages:
   alert: foobar
`)

	config.StdIn = bytes.NewBuffer(buf)

	msg, code := run("stdin", "INFO", false, false, "")

	assert.Equal(t, "error get core storages engine for alert 'foobar', storage not found", msg)
	assert.Equal(t, 1, code)
}

func Test_Main_error_init_alert_manager(t *testing.T) {
	buf := []byte(`
channels:
 syslog:
   - name: default
     network: tcp
     address: 127.0.0.1:70000
`)

	config.StdIn = bytes.NewBuffer(buf)

	msg, code := run("stdin", "INFO", false, false, "")

	assert.Equal(t, "error init alert manager, error init syslog channel default, dial tcp: address 70000: invalid port", msg)
	assert.Equal(t, 1, code)
}

func Test_Main_error_init_storages_kv(t *testing.T) {
	buf := []byte(`
global:
 storages:
   kv: foobar
`)

	config.StdIn = bytes.NewBuffer(buf)

	msg, code := run("stdin", "INFO", false, false, "")

	assert.Equal(t, "error get kv storage engine 'foobar', storage not found", msg)
	assert.Equal(t, 1, code)
}

func Test_validateLogLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "level ERROR",
			args: args{
				level: "ERROR",
			},
			wantErr: false,
		},
		{
			name: "level WARN",
			args: args{
				level: "WARN",
			},
			wantErr: false,
		},
		{
			name: "level INFO",
			args: args{
				level: "INFO",
			},
			wantErr: false,
		},
		{
			name: "level DEBUG",
			args: args{
				level: "DEBUG",
			},
			wantErr: false,
		},
		{
			name: "level wrong",
			args: args{
				level: "wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateLogLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("validateLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
