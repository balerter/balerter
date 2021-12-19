package http

import (
	"reflect"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func TestHTTP_parseRequestArgs(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		luaState *lua.LState
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
				a.luaState = lua.NewState()
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LNumber(10))
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LString("foo"))
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				URI:     "foo",
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LString("foo"))
				a.luaState.Push(lua.LNumber(42))
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LString("foo"))
				a.luaState.Push(lua.LString("bar"))
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				URI:     "foo",
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LString("foo"))
				a.luaState.Push(lua.LString("bar"))
				a.luaState.Push(lua.LString("baz"))
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
				a.luaState = lua.NewState()
				a.luaState.Push(lua.LString("foo"))
				a.luaState.Push(lua.LString("bar"))
				h := &lua.LTable{}
				h.RawSetString("Baz", lua.LString("Bar"))
				a.luaState.Push(h)
				return a
			}(),
			want: &requestArgs{
				Method:  "",
				URI:     "foo",
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
			}
			got, err := h.parseRequestArgs(tt.args.luaState)
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

func Test_parseMethod(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get",
			args: args{
				method: "get",
			},
			want: "GET",
		},
		{
			name: "bad",
			args: args{
				method: "bad",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMethod(tt.args.method); got != tt.want {
				t.Errorf("parseMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestArgs_parseFromTable(t *testing.T) {
	type fields struct {
		Method             string
		URI                string
		Body               []byte
		Headers            map[string]string
		InsecureSkipVerify bool
		Timeout            time.Duration
	}
	type args struct {
		tbl func() *lua.LTable
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		errValue string
	}{
		{
			name:   "empty",
			fields: fields{},
			args: args{
				tbl: func() *lua.LTable {
					return &lua.LTable{}
				},
			},
			wantErr:  false,
			errValue: "",
		},
		{
			name:   "bad method",
			fields: fields{},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("bad"))
					return tbl
				},
			},
			wantErr:  true,
			errValue: "bad http method bad",
		},
		{
			name:   "bad uri",
			fields: fields{},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("get"))
					tbl.RawSetString("uri", lua.LNumber(42))
					return tbl
				},
			},
			wantErr:  true,
			errValue: "uri must be a string",
		},
		{
			name:   "bad headers",
			fields: fields{},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("get"))
					tbl.RawSetString("uri", lua.LString("uri"))
					tbl.RawSetString("headers", lua.LNumber(42))
					return tbl
				},
			},
			wantErr:  true,
			errValue: "headers must be a table",
		},
		{
			name: "bad insecureSkipVerify",
			fields: fields{
				Method:  "GET",
				URI:     "uri",
				Body:    []byte("foo"),
				Headers: map[string]string{"a": "b"},
			},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("get"))
					tbl.RawSetString("uri", lua.LString("uri"))
					tbl.RawSetString("body", lua.LString("foo"))
					tbl.RawSetString("insecureSkipVerify", lua.LString("foo"))
					h := &lua.LTable{}
					h.RawSetString("a", lua.LString("b"))
					tbl.RawSetString("headers", h)
					return tbl
				},
			},
			wantErr:  true,
			errValue: "insecureSkipVerify must be a bool",
		},
		{
			name: "bad timeout",
			fields: fields{
				Method:  "GET",
				URI:     "uri",
				Body:    []byte("foo"),
				Headers: map[string]string{"a": "b"},
			},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("get"))
					tbl.RawSetString("uri", lua.LString("uri"))
					tbl.RawSetString("body", lua.LString("foo"))
					tbl.RawSetString("insecureSkipVerify", lua.LBool(true))
					tbl.RawSetString("timeout", lua.LString("10r"))
					h := &lua.LTable{}
					h.RawSetString("a", lua.LString("b"))
					tbl.RawSetString("headers", h)
					return tbl
				},
			},
			wantErr:  true,
			errValue: "timeout must be a time.Duration",
		},
		{
			name: "ok",
			fields: fields{
				Method:             "GET",
				URI:                "uri",
				Body:               []byte("foo"),
				Headers:            map[string]string{"a": "b"},
				InsecureSkipVerify: true,
				Timeout:            time.Second * 10,
			},
			args: args{
				tbl: func() *lua.LTable {
					tbl := &lua.LTable{}
					tbl.RawSetString("method", lua.LString("get"))
					tbl.RawSetString("uri", lua.LString("uri"))
					tbl.RawSetString("body", lua.LString("foo"))
					tbl.RawSetString("insecureSkipVerify", lua.LBool(true))
					tbl.RawSetString("timeout", lua.LString("10s"))
					h := &lua.LTable{}
					h.RawSetString("a", lua.LString("b"))
					tbl.RawSetString("headers", h)
					return tbl
				},
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestArgs{
				Method:  tt.fields.Method,
				URI:     tt.fields.URI,
				Body:    tt.fields.Body,
				Headers: tt.fields.Headers,
			}
			err := r.parseFromTable(tt.args.tbl())
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFromTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t.Errorf("unexpected error value = %s, want %s", err.Error(), tt.errValue)
			}
		})
	}
}
