package notify

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New(config.ChannelNotify{
		Name: "foo",
		Icons: config.ChannelNotifyIcons{
			Success: "icon1",
			Error:   "icon2",
			Warning: "icon3",
		},
	}, nil)

	require.NoError(t, err)
	assert.IsType(t, &Notify{}, p)
	assert.Equal(t, "foo", p.name)
	assert.Equal(t, "icon1", p.iconSuccess)
	assert.Equal(t, "icon2", p.iconError)
	assert.Equal(t, "icon3", p.iconWarning)
}

func TestName(t *testing.T) {
	p := &Notify{name: "foo"}
	assert.Equal(t, "foo", p.Name())
}

func TestNotify_getIconByLevel(t *testing.T) {
	type fields struct {
		iconSuccess string
		iconError   string
		iconWarning string
	}
	type args struct {
		l string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "success",
			fields: fields{iconSuccess: "icon1", iconError: "icon2", iconWarning: "icon3"},
			args:   args{l: "success"},
			want:   "icon1",
		},
		{
			name:   "warning",
			fields: fields{iconSuccess: "icon1", iconError: "icon2", iconWarning: "icon3"},
			args:   args{l: "warning"},
			want:   "icon3",
		},
		{
			name:   "error",
			fields: fields{iconSuccess: "icon1", iconError: "icon2", iconWarning: "icon3"},
			args:   args{l: "error"},
			want:   "icon2",
		},
		{
			name:   "unknown",
			fields: fields{iconSuccess: "icon1", iconError: "icon2", iconWarning: "icon3"},
			args:   args{l: "unknown"},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Notify{
				iconSuccess: tt.fields.iconSuccess,
				iconError:   tt.fields.iconError,
				iconWarning: tt.fields.iconWarning,
			}
			if got := p.getIconByLevel(tt.args.l); got != tt.want {
				t.Errorf("getIconByLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
