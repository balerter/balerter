package alerts

import (
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_parseNames(t *testing.T) {
	type args struct {
		argValue string
	}
	tests := []struct {
		name string
		args args
		want map[string]struct{}
	}{
		{
			name: "empty argument",
			args: args{
				argValue: "",
			},
			want: map[string]struct{}{},
		},
		{
			name: "one item",
			args: args{
				argValue: "name1",
			},
			want: map[string]struct{}{"name1": {}},
		},
		{
			name: "two items",
			args: args{
				argValue: "name1,name2",
			},
			want: map[string]struct{}{"name1": {}, "name2": {}},
		},
		{
			name: "with empty",
			args: args{
				argValue: "name1,,name2",
			},
			want: map[string]struct{}{"name1": {}, "": {}, "name2": {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNames(tt.args.argValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLevels(t *testing.T) {
	type args struct {
		argValue string
	}
	tests := []struct {
		name    string
		args    args
		want    map[alert2.Level]struct{}
		wantErr bool
		errText string
	}{
		{
			name: "empty arg",
			args: args{
				argValue: "",
			},
			want:    map[alert2.Level]struct{}{},
			wantErr: false,
			errText: "",
		},
		{
			name: "one item",
			args: args{
				argValue: "error",
			},
			want:    map[alert2.Level]struct{}{alert2.LevelError: {}},
			wantErr: false,
			errText: "",
		},
		{
			name: "two items",
			args: args{
				argValue: "error,success",
			},
			want:    map[alert2.Level]struct{}{alert2.LevelError: {}, alert2.LevelSuccess: {}},
			wantErr: false,
			errText: "",
		},
		{
			name: "with bad value: empty",
			args: args{
				argValue: "error,,success",
			},
			want:    nil,
			wantErr: true,
			errText: "bad level value",
		},
		{
			name: "with bad value: wrong",
			args: args{
				argValue: "error,foo,success",
			},
			want:    nil,
			wantErr: true,
			errText: "bad level value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevels(tt.args.argValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLevels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errText {
				t.Errorf("parseLevels() error = '%s', wantErrText '%s'", err.Error(), tt.errText)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLevels() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_ErrorLevel(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{},
	}
	req.URL.RawQuery = "name=foo,bar&level=error,success,bad"

	var data []*alert2.Alert

	_, err := filter(req, data)
	require.Error(t, err)
	assert.Equal(t, "bad level value", err.Error())
}

func TestFilter(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{},
	}
	req.URL.RawQuery = "name=foo,bar&level=error,success"

	var data []*alert2.Alert

	a0 := alert2.AcquireAlert()
	a0.UpdateLevel(alert2.LevelSuccess)
	a0.SetName("baz")
	data = append(data, a0)

	a1 := alert2.AcquireAlert()
	a1.UpdateLevel(alert2.LevelSuccess)
	a1.SetName("foo")
	data = append(data, a1)

	a2 := alert2.AcquireAlert()
	a2.UpdateLevel(alert2.LevelError)
	a2.SetName("bar")
	data = append(data, a2)

	a3 := alert2.AcquireAlert()
	a3.UpdateLevel(alert2.LevelWarn)
	a3.SetName("bar")
	data = append(data, a3)

	result, err := filter(req, data)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result))
}
