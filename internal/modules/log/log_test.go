package log

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestLog_Loader(t *testing.T) {
	logger := &Log{}

	L := lua.NewState()

	n := logger.Loader(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	assert.IsType(t, &lua.LNilType{}, v.RawGet(lua.LString("wrong-name")))

	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("error")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("warn")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("info")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("debug")))
}

func TestLog_levels(t *testing.T) {

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	lg := &Log{
		logger:  logger,
		jobName: "jobname",
	}

	L := lua.NewState()
	L.Push(lua.LString("message error"))
	lg.error(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message warn").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message info").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("job", lg.jobName)).Len())

	L = lua.NewState()
	L.Push(lua.LString("message warn"))
	lg.warn(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message info").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("job", lg.jobName)).Len())

	L = lua.NewState()
	L.Push(lua.LString("message info"))
	lg.info(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message info").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("job", lg.jobName)).Len())

	L = lua.NewState()
	L.Push(lua.LString("message debug"))
	lg.debug(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message info").FilterField(zap.String("job", lg.jobName)).Len())
	assert.Equal(t, 1, logs.FilterMessage("message debug").FilterField(zap.String("job", lg.jobName)).Len())
}

func TestNew(t *testing.T) {
	lgfunc := New("jobname", zap.NewNop())
	assert.IsType(t, func() lua.LGFunction { return func(L *lua.LState) int { return 0 } }(), lgfunc)
}
