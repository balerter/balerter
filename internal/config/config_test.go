package config

import (
	"flag"
	"fmt"
	"github.com/balerter/balerter/internal/config/api"
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/balerter/balerter/internal/config/datasources"
	clickhouseCfg "github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/config/scripts"
	"github.com/balerter/balerter/internal/config/storages/core"
	"github.com/balerter/balerter/internal/config/storages/core/sqlite"
	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/balerter/balerter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_checkUnique(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty slice",
			args: args{
				data: []string{},
			},
			want: "",
		},
		{
			name: "one element",
			args: args{
				data: []string{"foo"},
			},
			want: "",
		},
		{
			name: "without duplicates",
			args: args{
				data: []string{"foo", "bar"},
			},
			want: "",
		},
		{
			name: "with duplicates",
			args: args{
				data: []string{"foo", "bar", "baz", "foo"},
			},
			want: "foo",
		},
		{
			name: "with duplicates case insensitive",
			args: args{
				data: []string{"foo", "bar", "baz", "FOO"},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.CheckUnique(tt.args.data); got != tt.want {
				t.Errorf("checkUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeCfg_unknown_format(t *testing.T) {
	err := decodeCfg("foo.bar", nil, nil)
	require.Error(t, err)
	assert.Equal(t, "unknown format", err.Error())
}

func Test_decodeCfg_yml(t *testing.T) {
	cfg := &Config{}
	err := decodeCfg("foo.yml", []byte(`storageKV: foo`), cfg)
	require.NoError(t, err)
	assert.Equal(t, "foo", cfg.StorageKV)
}

func Test_decodeCfg_hcl(t *testing.T) {
	cfg := &Config{}
	err := decodeCfg("foo.hcl", []byte(`storageKV = "foo"`), cfg)
	require.NoError(t, err)
	assert.Equal(t, "foo", cfg.StorageKV)
}

func TestNew_error_read_file(t *testing.T) {
	_, _, err := New(flag.NewFlagSet("fs", flag.ContinueOnError), []string{"-config=notfound"})
	require.Error(t, err)
	assert.Equal(t, "error read config file, open notfound: no such file or directory", err.Error())
}

type badReader struct{}

func (br *badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("err1")
}

func TestNew_error_read_stdin(t *testing.T) {
	var err error
	StdIn = &badReader{}
	_, _, err = New(flag.NewFlagSet("fs", flag.ContinueOnError), []string{"-config=stdin"})
	require.Error(t, err)
	assert.Equal(t, "error read config file, err1", err.Error())
}

func TestNew_error_decode(t *testing.T) {
	f, err := tempConfigFile(".bad")
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	_, _, err = New(flag.NewFlagSet("fs", flag.ContinueOnError), []string{"-config=" + f.Name()})
	require.Error(t, err)
	assert.Equal(t, "error parse config file, unknown format", err.Error())
}

func TestNew_error_validation(t *testing.T) {
	c := `
storagesCore {
  sqlite "filestorage1" {
    path = "/tmp/sqlite.db"
    tableKV {
      table = ""
      fields {
        key = "key"
        value = "value"
      }
    }
    tableAlerts {
      table = ""
      fields {
        name = "id"
        level = "level"
        count = "count"
        createdAt = "created_at"
        updatedAt = "updated_at"
      }
    }
  }
}
`

	f, err := tempConfigFile(".hcl")
	require.NoError(t, err)
	_, err = f.Write([]byte(c))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	_, _, err = New(flag.NewFlagSet("fs", flag.ContinueOnError), []string{"-config=" + f.Name()})
	require.Error(t, err)
	assert.Equal(t, "error config validation, error StoragesCore validation, table must be not empty", err.Error())
}

func TestNew_error_parse_flags(t *testing.T) {
	_, _, err := New(
		flag.NewFlagSet("fs", flag.ContinueOnError),
		[]string{"-foo"},
	)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "flag provided but not defined: -foo"))
}

func TestNew(t *testing.T) {
	c := `
storagesCore {
  sqlite "filestorage1" {
    path = "/tmp/sqlite.db"
    tableKV {
      table = "1"
      fields {
        key = "key"
        value = "value"
      }
    }
    tableAlerts {
      table = "2"
      fields {
        name = "id"
        level = "level"
        count = "count"
        createdAt = "created_at"
        updatedAt = "updated_at"
      }
    }
  }
}
`

	f, err := tempConfigFile(".hcl")
	require.NoError(t, err)
	_, err = f.Write([]byte(c))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	cfg, flags, err := New(
		flag.NewFlagSet("fs", flag.ContinueOnError),
		[]string{"-config=" + f.Name(), "-once", "-debug"},
	)
	require.NoError(t, err)
	assert.Equal(t, "filestorage1", cfg.StoragesCore.Sqlite[0].Name)
	assert.Equal(t, f.Name(), flags.ConfigFilePath)
	assert.True(t, flags.Once)
	assert.True(t, flags.Debug)
}

func tempConfigFile(ext string) (*os.File, error) {
	rand.Seed(time.Now().UnixNano())
	return os.Create(path.Join(os.TempDir(), strconv.Itoa(1e6)+ext))
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Scripts        *scripts.Scripts
		DataSources    *datasources.DataSources
		Channels       *channels.Channels
		StoragesUpload *upload.Upload
		StoragesCore   *core.Core
		API            *api.API
		LuaModulesPath string
		StorageAlert   string
		StorageKV      string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name: "bad scripts",
			fields: fields{
				Scripts: &scripts.Scripts{
					UpdateInterval: -1,
					Folder:         nil,
					File:           nil,
					Postgres:       nil,
				},
			},
			wantErr:  true,
			errValue: "error Scripts validation, updateInterval must be not less than 0",
		},
		{
			name: "bad datasource",
			fields: fields{
				DataSources: &datasources.DataSources{
					Clickhouse: []clickhouseCfg.Clickhouse{{}},
				},
			},
			wantErr:  true,
			errValue: "error DataSources validation, name must be not empty",
		},
		{
			name: "bad channels",
			fields: fields{
				Channels: &channels.Channels{
					Email: []email.Email{{}},
				},
			},
			wantErr:  true,
			errValue: "error Channels validation, validate channel email: name must be not empty",
		},
		{
			name: "bad storagesUpload",
			fields: fields{
				StoragesUpload: &upload.Upload{
					S3: []s3.S3{{}},
				},
			},
			wantErr:  true,
			errValue: "error StoragesUpload validation, name must be not empty",
		},
		{
			name: "bad storagesCore",
			fields: fields{
				StoragesCore: &core.Core{
					Sqlite: []sqlite.Sqlite{{}},
				},
			},
			wantErr:  true,
			errValue: "error StoragesCore validation, name must be not empty",
		},
		{
			name: "ok, with api",
			fields: fields{
				API: &api.API{},
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Scripts:        tt.fields.Scripts,
				DataSources:    tt.fields.DataSources,
				Channels:       tt.fields.Channels,
				StoragesUpload: tt.fields.StoragesUpload,
				StoragesCore:   tt.fields.StoragesCore,
				API:            tt.fields.API,
				LuaModulesPath: tt.fields.LuaModulesPath,
				StorageAlert:   tt.fields.StorageAlert,
				StorageKV:      tt.fields.StorageKV,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t.Errorf("unexpected error value = %s, want = %s", err.Error(), tt.errValue)
			}
		})
	}
}
