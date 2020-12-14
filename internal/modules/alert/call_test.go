package alert

import (
	alert2 "github.com/balerter/balerter/internal/alert"
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
		},
		{
			name:   "with options",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))

					opts := &lua.LTable{}
					opts.RawSet(lua.LString("quiet"), lua.LBool(true))
					fields := &lua.LTable{}
					fields.RawSetInt(1, lua.LString("foo"))
					fields.RawSetInt(2, lua.LString("bar"))
					opts.RawSet(lua.LString("fields"), fields)

					L.Push(opts)

					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Quiet: true, Fields: []string{"foo", "bar"}},
			wantErr:          false,
		},
		{
			name:   "with wrong options",
			fields: defaultFields,
			args: args{
				luaState: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))

					opts := &lua.LTable{}
					opts.RawSet(lua.LString("quiet"), lua.LString("not bool"))
					fields := &lua.LTable{}
					fields.RawSetInt(1, lua.LString("foo"))
					fields.RawSetInt(2, lua.LString("bar"))
					opts.RawSet(lua.LString("fields"), fields)

					L.Push(opts)

					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: &alert2.Options{Quiet: false, Fields: []string{"foo", "bar"}},
			wantErr:          true,
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

func TestManager_luaCall_errorGetAlertData(t *testing.T) {
	m := &Alert{
		logger: zap.NewNop(),
	}

	opts := &lua.LTable{}
	opts.RawSet(lua.LString("repeat"), lua.LString("wrong value"))

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	L.Push(opts)
	n := m.call(script.New(), alert2.LevelError)(L)
	assert.Equal(t, 1, n)

	v := L.Get(4).String()

	assert.Equal(t, "error get arguments: wrong options format: 1 error(s) decoding:\n\n* cannot "+
		"parse 'Repeat' as int: strconv.ParseInt: parsing \"wrong value\": invalid syntax", v)
}
