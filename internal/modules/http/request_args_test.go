package http

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"testing"
)

func TestHTTP_parseRequestArgs(t *testing.T) {
	type fields struct {
		logger *zap.Logger
		client *http.Client
	}
	type args struct {
		L *lua.LState
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *requestArgs
		wantErr   bool
		errorText string
	}{
		{
			name:   "empty fields",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				return a
			}(),
			want:      nil,
			wantErr:   true,
			errorText: "uri argument must be a string",
		},
		{
			name:   "uri is number",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LNumber(10))
				return a
			}(),
			want:      nil,
			wantErr:   true,
			errorText: "uri argument must be a string",
		},
		{
			name:   "uri is string",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LString("foo"))
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				Uri:     "foo",
				Body:    nil,
				Headers: map[string]string{},
			},
			wantErr:   false,
			errorText: "",
		},
		{
			name:   "body is not a string",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LString("foo"))
				a.L.Push(lua.LNumber(42))
				return a
			}(),
			want:      nil,
			wantErr:   true,
			errorText: "body argument must be a string or nil",
		},
		{
			name:   "uri and body",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LString("foo"))
				a.L.Push(lua.LString("bar"))
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				Uri:     "foo",
				Body:    []byte("bar"),
				Headers: map[string]string{},
			},
			wantErr:   false,
			errorText: "",
		},
		{
			name:   "headers is not a table",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LString("foo"))
				a.L.Push(lua.LString("bar"))
				a.L.Push(lua.LString("baz"))
				return a
			}(),
			want:      nil,
			wantErr:   true,
			errorText: "headers argument must be a table or nil",
		},
		{
			name:   "uri, body and headers",
			fields: fields{},
			args: func() args {
				a := args{}
				a.L = lua.NewState()
				a.L.Push(lua.LString("foo"))
				a.L.Push(lua.LString("bar"))
				h := &lua.LTable{}
				h.RawSetString("Baz", lua.LString("Bar"))
				a.L.Push(h)
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				Uri:     "foo",
				Body:    []byte("bar"),
				Headers: map[string]string{"Baz": "Bar"},
			},
			wantErr:   false,
			errorText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTP{
				logger: tt.fields.logger,
				client: tt.fields.client,
			}
			got, err := h.parseRequestArgs(tt.args.L)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRequestArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.errorText != "" && err.Error() != tt.errorText {
				t.Errorf("parseRequestArgs() error = %s, errorText %s", err.Error(), tt.errorText)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRequestArgs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
