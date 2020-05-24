package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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

func (m *alertChannelMock) Send(mes *message.Message) error {
	args := m.Called(mes)
	return args.Error(0)
}

func TestManager_getAlertData(t *testing.T) {
	type fields struct {
		logger   *zap.Logger
		channels map[string]alertChannel
	}

	type args struct {
		luaState *lua.LState
	}

	defaultFields := fields{
		logger:   zap.NewNop(),
		channels: map[string]alertChannel{},
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
				luaState: func() *lua.LState {
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
				luaState: func() *lua.LState {
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
				luaState: func() *lua.LState {
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
				luaState: func() *lua.LState {
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
			wantAlertOptions: options{},
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
			wantAlertOptions: options{Quiet: true, Fields: []string{"foo", "bar"}},
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
			wantAlertOptions: options{Quiet: false, Fields: []string{"foo", "bar"}},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				logger:   tt.fields.logger,
				channels: tt.fields.channels,
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
	chan1 := &alertChannelMock{}

	m := &Manager{
		logger:   zap.NewNop(),
		channels: map[string]alertChannel{"chan1": chan1},
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

	assert.Equal(t, "error get arguments: wrong options format: 1 error(s) decoding:\n\n* cannot "+
		"parse 'Repeat' as int: strconv.ParseInt: parsing \"wrong value\": invalid syntax", v)
}

func TestManager_luaCall_error_get_alert(t *testing.T) {
	chan1 := &alertChannelMock{}
	eng := coreStorage.NewMock("")

	eng.AlertMock().On("GetOrNew", mock.Anything).Return(nil, fmt.Errorf("error1"))

	m := &Manager{
		logger:   zap.NewNop(),
		engine:   eng,
		channels: map[string]alertChannel{"chan1": chan1},
	}

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	n := m.luaCall(script.New(), alert.LevelError)(L)
	assert.Equal(t, 1, n)

	v := L.Get(3).String()

	assert.Equal(t, "internal error get alert from storage: error1", v)
}

func TestManager_luaCall_change_level(t *testing.T) {
	chan1 := &alertChannelMock{}
	eng := coreStorage.NewMock("")

	a := alert.AcquireAlert()

	chan1.On("Send", mock.Anything).Run(func(args mock.Arguments) {
		m := args.Get(0).(*message.Message)

		assert.Equal(t, "alertName", m.AlertName)
		assert.Equal(t, "alertText1", m.Text)
		assert.Equal(t, alert.LevelError.String(), m.Level)
	}).Return(nil)

	eng.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	eng.AlertMock().On("Release", mock.Anything).Return(nil)

	m := &Manager{
		logger:   zap.NewNop(),
		engine:   eng,
		channels: map[string]alertChannel{"chan1": chan1},
	}

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	n := m.luaCall(script.New(), alert.LevelError)(L)
	assert.Equal(t, 0, n)

	assert.Equal(t, alert.LevelError, a.Level())

	chan1.AssertCalled(t, "Send", mock.Anything)
	chan1.AssertExpectations(t)
}

func TestManager_luaCall_same_level(t *testing.T) {
	chan1 := &alertChannelMock{}
	eng := coreStorage.NewMock("")

	a := alert.AcquireAlert()
	a.UpdateLevel(alert.LevelError)

	eng.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	eng.AlertMock().On("Release", mock.Anything).Return(nil)
	chan1.On("Send", mock.Anything).Run(func(args mock.Arguments) {
		m := args.Get(0).(*message.Message)

		assert.Equal(t, "alertName", m.AlertName)
		assert.Equal(t, "alertText1", m.Text)
		assert.Equal(t, alert.LevelError.String(), m.Level)
	}).Return(nil)

	m := &Manager{
		logger:   zap.NewNop(),
		engine:   eng,
		channels: map[string]alertChannel{"chan1": chan1},
	}

	opts := &lua.LTable{}
	opts.RawSet(lua.LString("repeat"), lua.LNumber(1))

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	L.Push(opts)
	n := m.luaCall(script.New(), alert.LevelError)(L)
	assert.Equal(t, 0, n)

	assert.Equal(t, alert.LevelError, a.Level())
	assert.Equal(t, 1, a.Count())

	chan1.AssertCalled(t, "Send", mock.Anything)
	chan1.AssertExpectations(t)
}

func TestManager_luaCall_error_release(t *testing.T) {
	chan1 := &alertChannelMock{}
	eng := coreStorage.NewMock("")

	a := alert.AcquireAlert()

	chan1.On("Send", mock.Anything).Return(nil)
	eng.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	eng.AlertMock().On("Release", mock.Anything).Return(fmt.Errorf("error1"))

	core, logs := observer.New(zap.ErrorLevel)

	m := &Manager{
		logger:   zap.New(core),
		engine:   eng,
		channels: map[string]alertChannel{"chan1": chan1},
	}

	L := lua.NewState()
	L.Push(lua.LString("alertName"))
	L.Push(lua.LString("alertText1"))
	n := m.luaCall(script.New(), alert.LevelError)(L)
	assert.Equal(t, 0, n)

	assert.Equal(t, alert.LevelError, a.Level())

	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, 1, logs.FilterField(zap.Error(fmt.Errorf("error1"))).FilterMessage("error release alert").Len())
}
