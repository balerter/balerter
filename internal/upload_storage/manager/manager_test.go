package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type moduleMock struct {
	name string
	mock.Mock
}

func (m *moduleMock) Name() string {
	return ""
}

func (m *moduleMock) GetLoader(_ *script.Script) lua.LGFunction {
	return nil
}

func (m *moduleMock) Stop() error {
	return nil
}

func TestNew(t *testing.T) {
	m := New(zap.NewNop())
	assert.IsType(t, &Manager{}, m)
}

func TestManager_Get(t *testing.T) {
	type fields struct {
		logger  *zap.Logger
		modules map[string]modules.Module
	}

	m1 := &moduleMock{name: "foo"}
	m2 := &moduleMock{name: "bar"}

	tests := []struct {
		name   string
		fields fields
		want   []modules.Module
	}{
		{
			name: "",
			fields: fields{
				logger: nil,
				modules: map[string]modules.Module{
					"foo": m1,
					"bar": m2,
				},
			},
			want: []modules.Module{m1, m2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				logger:  tt.fields.logger,
				modules: tt.fields.modules,
			}
			if got := m.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Init(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.Module{},
	}

	err := m.Init(config.StoragesUpload{S3: []config.StorageUploadS3{{
		Name:     "f1",
		Region:   "f2",
		Key:      "f3",
		Secret:   "f4",
		Endpoint: "f5",
		Bucket:   "f6",
	}}})

	require.NoError(t, err)

	md, ok := m.modules["s3.f1"]
	require.True(t, ok)

	assert.Equal(t, "s3.f1", md.Name())
}
