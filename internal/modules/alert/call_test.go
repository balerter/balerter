package alert

import (
	alert2 "github.com/balerter/balerter/internal/alert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestManager_getAlertData(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}

	type args struct {
		luaState *lua.LState
	}

	defaultFields := fields{
		logger: zap.NewNop(),
	}

	tests := []struct {
		name             string
		fields           fields
		args             args
		wantAlertName    string
		wantAlertText    string
		wantAlertOptions *alert2.Options
		wantErr          bool
		wantErrString    string
	}{
		{
			name:   "empty args",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					return L
				}(),
			},
			wantAlertName:    "",
			wantAlertText:    "",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "alert name must be provided",
		},
		{
			name:   "only alert name",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "",
			wantAlertOptions: &alert2.Options{},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "empty (only space) alert name",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString(" "))
					return L
				}(),
			},
			wantAlertName:    "",
			wantAlertText:    "",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "alert name must be not empty",
		},
		{
			name:   "alert name and text",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "with options NOT table",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					L.Push(lua.LString("options"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "options must be a table",
		},
		{
			name:   "with options - bad channels, not a table",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("channels"), lua.LBool(true))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "channels options must be a table",
		},
		{
			name:   "with options - bad channels, not string value",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					channels := &lua.LTable{}
					channels.RawSetInt(1, lua.LString("foo"))
					channels.RawSetInt(2, lua.LNumber(2))
					opts.RawSet(lua.LString("channels"), channels)
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "wrong channel name 2",
		},
		{
			name:   "with options - channels",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					channels := &lua.LTable{}
					channels.RawSetInt(1, lua.LString("foo"))
					channels.RawSetInt(2, lua.LString("bar"))
					opts.RawSet(lua.LString("channels"), channels)
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Channels: []string{"foo", "bar"}},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "with options - bad quiet",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("quiet"), lua.LNumber(2))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "quiet must be a bool",
		},
		{
			name:   "with options - quiet",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("quiet"), lua.LBool(true))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Quiet: true},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "with options - repeat, not a number",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("repeat"), lua.LString("foo"))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "repeat must be a number",
		},
		{
			name:   "with options - repeat",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("repeat"), lua.LNumber(42))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Repeat: 42},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "with not int resend",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("resend"), lua.LString("not int"))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "resend must be a number",
		},
		{
			name:   "with resend and repeat",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("resend"), lua.LNumber(1))
					opts.RawSet(lua.LString("repeat"), lua.LNumber(2))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "you must not use repeat and resend option together",
		},
		{
			name:   "with resend",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("resend"), lua.LNumber(42))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Repeat: 42},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "with image not a string",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("image"), lua.LNumber(42))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{},
			wantErr:          true,
			wantErrString:    "image must be a string",
		},
		{
			name:   "with image",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("image"), lua.LString("foo"))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Image: "foo"},
			wantErr:          false,
			wantErrString:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Alert{
				logger: tt.fields.logger,
			}
			gotAlertName, gotAlertText, gotAlertOptions, err := m.getAlertData(tt.args.luaState)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAlertData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.wantErrString != err.Error() {
					t.Errorf("error string = '%v', want '%v'", err.Error(), tt.wantErrString)
				}
				return
			}
			if gotAlertName != tt.wantAlertName {
				t.Errorf("getAlertData() gotAlertName = %v, want %v", gotAlertName, tt.wantAlertName)
			}
			if gotAlertText != tt.wantAlertText {
				t.Errorf("getAlertData() gotAlertText = %v, want %v", gotAlertText, tt.wantAlertText)
			}
			if !reflect.DeepEqual(gotAlertOptions, tt.wantAlertOptions) {
				t.Errorf("getAlertData() gotAlertOptions = %v, want %v", gotAlertOptions, tt.wantAlertOptions)
			}
		})
	}
}
