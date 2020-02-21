package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type alertChannelMock struct {
	mock.Mock
}

func (m *alertChannelMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *alertChannelMock) Send(message *message.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestManager_getAlertData(t *testing.T) {
	type fields struct {
		logger   *zap.Logger
		channels map[string]alertChannel
		alerts   map[string]*alert.Alert
	}

	type args struct {
		L *lua.LState
	}

	defaultFields := fields{
		logger:   zap.NewNop(),
		channels: map[string]alertChannel{},
		alerts:   map[string]*alert.Alert{},
	}

	tests := []struct {
		name             string
		fields           fields
		args             args
		wantAlertName    string
		wantAlertText    string
		wantAlertOptions options
		wantErr          bool
	}{
		{
			name:   "empty args",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
					L := lua.NewState()
					return L
				}(),
			},
			wantAlertName:    "",
			wantAlertText:    "",
			wantAlertOptions: options{},
			wantErr:          true,
		},
		{
			name:   "only alert name",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "",
			wantAlertOptions: options{},
			wantErr:          false,
		},
		{
			name:   "empty (only space) alert name",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString(" "))
					return L
				}(),
			},
			wantAlertName:    "",
			wantAlertText:    "",
			wantAlertOptions: options{},
			wantErr:          true,
		},
		{
			name:   "alert name and text",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: options{},
			wantErr:          false,
		},
		{
			name:   "with options NOT table",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
					L := lua.NewState()
					L.Push(lua.LString("alertName1"))
					L.Push(lua.LString("alertText1"))
					L.Push(lua.LString("options"))
					return L
				}(),
			},
			wantAlertName:    "alertName1",
			wantAlertText:    "alertText1",
			wantAlertOptions: options{},
			wantErr:          true,
		},
		{
			name:   "with options",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
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
			wantAlertOptions: options{Quiet: true, Fields: []string{"foo", "bar"}},
			wantErr:          false,
		},
		{
			name:   "with wrong options",
			fields: defaultFields,
			args: args{
				L: func() *lua.LState {
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
			wantAlertOptions: options{Quiet: false, Fields: []string{"foo", "bar"}},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				logger:   tt.fields.logger,
				channels: tt.fields.channels,
				alerts:   tt.fields.alerts,
			}
			gotAlertName, gotAlertText, gotAlertOptions, err := m.getAlertData(tt.args.L)
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
	chan1 := &alertChannelMock{}

	m := &Manager{
		logger:   zap.NewNop(),
		channels: map[string]alertChannel{"chan1": chan1},
		alerts:   map[string]*alert.Alert{},
	}

	opts := &lua.LTable{}
	opts.RawSet(lua.LString("repeat"), lua.LString("wrong value"))

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	L.Push(opts)
	n := m.luaCall(script.New(), alert.LevelError)(L)
	assert.Equal(t, 1, n)

	v := L.Get(4).String()

	assert.Equal(t, "error get arguments: wrong options format: 1 error(s) decoding:\n\n* cannot parse 'Repeat' as int: strconv.ParseInt: parsing \"wrong value\": invalid syntax", v)
}

// todo fix test - without use `&message.Message`
//func TestManager_luaCall_repeat(t *testing.T) {
//	chan1 := &alertChannelMock{}
//	chan1.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)
//
//	m := &Manager{
//		logger:   zap.NewNop(),
//		channels: map[string]alertChannel{"chan1": chan1},
//		alerts:   map[string]*alert.Alert{},
//	}
//
//	opts := &lua.LTable{}
//	opts.RawSet(lua.LString("repeat"), lua.LNumber(2))
//
//	L := lua.NewState()
//	L.Push(lua.LString("alertName"))
//	L.Push(lua.LString("alertText1"))
//	L.Push(opts)
//	n := m.luaCall(script.New(), alert.LevelError)(L)
//	assert.Equal(t, 0, n)
//	chan1.AssertCalled(t, "Send", alert.LevelError, &message.Message{AlertName: "alertName", Text: "alertText1"}, nil)
//
//	L = lua.NewState()
//	L.Push(lua.LString("alertName"))
//	L.Push(lua.LString("alertText2"))
//	L.Push(opts)
//	n = m.luaCall(script.New(), alert.LevelError)(L)
//	assert.Equal(t, 0, n)
//	chan1.AssertNotCalled(t, "Send", alert.LevelError, &message.Message{AlertName: "alertName", Text: "alertText2"}, nil)
//
//	L = lua.NewState()
//	L.Push(lua.LString("alertName"))
//	L.Push(lua.LString("alertText3"))
//	L.Push(opts)
//	n = m.luaCall(script.New(), alert.LevelError)(L)
//	assert.Equal(t, 0, n)
//	chan1.AssertCalled(t, "Send", alert.LevelError, &message.Message{AlertName: "alertName", Text: "alertText3"}, nil)
//
//	chan1.AssertExpectations(t)
//}
