package alert

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
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
		{
			name:   "fields ok",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					fields := &lua.LTable{}
					fields.RawSetString("foo", lua.LString("bar"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("fields"), fields)
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Fields: map[string]string{"foo": "bar"}},
			wantErr:          false,
			wantErrString:    "",
		},
		{
			name:   "fields wrong key type",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					fields := &lua.LTable{}
					fields.RawSet(lua.LNumber(42), lua.LString("bar"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("fields"), fields)
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Fields: map[string]string{"foo": "bar"}},
			wantErr:          true,
			wantErrString:    "option key must be a string, 42",
		},
		{
			name:   "fields wrong value type",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					fields := &lua.LTable{}
					fields.RawSetString("foo", lua.LNumber(42))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("fields"), fields)
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Fields: map[string]string{"foo": "bar"}},
			wantErr:          true,
			wantErrString:    "option value must be a string, 42",
		},
		{
			name:   "fields not a table",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					opts := &lua.LTable{}
					opts.RawSet(lua.LString("fields"), lua.LNumber(42))
					L.Push(opts)
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Fields: map[string]string{"foo": "bar"}},
			wantErr:          true,
			wantErrString:    "fields option must be a table",
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
			if !optionsAreEqual(*gotAlertOptions, *tt.wantAlertOptions) {
				t.Errorf("getAlertData() gotAlertOptions = %v, want %v", gotAlertOptions, tt.wantAlertOptions)
			}
		})
	}
}

func optionsAreEqual(got, want alert2.Options) bool {
	if !reflect.DeepEqual(got.Channels, want.Channels) {
		return false
	}
	if got.Quiet != want.Quiet {
		return false
	}
	if got.Repeat != want.Repeat {
		return false
	}
	if got.Image != want.Image {
		return false
	}
	for k, v := range got.Fields {
		wantV, ok := want.Fields[k]
		if !ok {
			return false
		}
		if wantV != v {
			return false
		}
		delete(want.Fields, k)
	}

	if len(want.Fields) > 0 {
		return false
	}

	return true
}

func TestAlert_call_error_get_alertData(t *testing.T) {
	a := &Alert{
		logger: zap.NewNop(),
	}

	f := a.callFromLua(nil, map[int][]string{}, alert2.LevelError)

	ls := lua.NewState()

	n := f(ls)
	assert.Equal(t, 1, n)
	assert.Equal(t, "error get arguments: alert name must be provided", ls.Get(1).String())
}

func TestAlert_call_error_update(t *testing.T) {
	am := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return nil, false, fmt.Errorf("err1")
		},
	}

	a := &Alert{
		logger:  zap.NewNop(),
		storage: am,
	}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := a.callFromLua(j.Script().Channels, map[int][]string{}, alert2.LevelError)

	ls := lua.NewState()
	ls.Push(lua.LString("foo"))

	n := f(ls)
	assert.Equal(t, 1, n)
	assert.Equal(t, "error update an alert: err1", ls.Get(2).String())
}

func TestAlert_call_level_was_updated(t *testing.T) {
	ra := &alert2.Alert{}

	am := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return ra, true, nil
		},
	}

	chManager := &chManagerMock{
		SendFunc: func(_ *alert2.Alert, _ string, _ *alert2.Options) {

		},
	}

	a := &Alert{
		logger:    zap.NewNop(),
		storage:   am,
		chManager: chManager,
	}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := a.callFromLua(j.Script().Channels, map[int][]string{}, alert2.LevelError)

	ls := lua.NewState()
	ls.Push(lua.LString("foo"))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(chManager.SendCalls()))
}

func TestAlert_call_level_was_not_updated(t *testing.T) {
	ra := &alert2.Alert{}

	am := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return ra, false, nil
		},
	}

	chManager := &chManagerMock{
		SendFunc: func(_ *alert2.Alert, _ string, _ *alert2.Options) {

		},
	}

	a := &Alert{
		logger:    zap.NewNop(),
		storage:   am,
		chManager: chManager,
	}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := a.callFromLua(j.Script().Channels, map[int][]string{}, alert2.LevelError)

	ls := lua.NewState()
	ls.Push(lua.LString("foo"))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, len(chManager.SendCalls()))
}

func TestAlert_call_repeat(t *testing.T) {
	ra := &alert2.Alert{
		Count: 10,
	}

	am := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return ra, false, nil
		},
	}

	chManager := &chManagerMock{
		SendFunc: func(_ *alert2.Alert, _ string, _ *alert2.Options) {

		},
	}

	a := &Alert{
		logger:    zap.NewNop(),
		storage:   am,
		chManager: chManager,
	}

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
	}

	f := a.callFromLua(j.Script().Channels, map[int][]string{}, alert2.LevelError)

	ls := lua.NewState()
	ls.Push(lua.LString("id"))
	ls.Push(lua.LString("text"))
	opts := &lua.LTable{}
	opts.RawSetString("resend", lua.LNumber(2))
	ls.Push(opts)

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(chManager.SendCalls()))
}

func TestAlert_call_escalate(t *testing.T) {
	alrt := &alert2.Alert{
		Count: 10,
		Level: alert2.LevelError,
	}
	moduleAlertMock := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return alrt, false, nil
		},
	}

	var sentToChannels []string

	chManagerMock := &chManagerMock{
		SendFunc: func(_ *alert2.Alert, _ string, opts *alert2.Options) {
			sentToChannels = append(sentToChannels, opts.Channels...)
		},
	}

	a := &Alert{
		storage:   moduleAlertMock,
		chManager: chManagerMock,
		logger:    zap.NewNop(),
	}

	f := a.callFromLua(nil, map[int][]string{10: {"foo", "bar"}}, alert2.LevelError)

	ls := lua.NewState()
	ls.Push(lua.LString("id"))
	ls.Push(lua.LString("text"))

	n := f(ls)

	assert.Equal(t, 0, n)
	assert.Equal(t, []string{"foo", "bar"}, sentToChannels)
}
