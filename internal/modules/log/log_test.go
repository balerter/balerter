package log

import (
	"testing"

	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLog_Loader(t *testing.T) {
	logger := &Log{}

	L := lua.NewState()

	j := &modules.JobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{Name: "scriptName"}
		},
	}

	f := logger.GetLoader(j)
	n := f(L)
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
		logger: logger,
	}

	L := lua.NewState()
	L.Push(lua.LString("message error"))
	lg.error("scriptName")(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message warn").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message info").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("scriptName", "scriptName")).Len())

	L = lua.NewState()
	L.Push(lua.LString("message warn"))
	lg.warn("scriptName")(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message info").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("scriptName", "scriptName")).Len())

	L = lua.NewState()
	L.Push(lua.LString("message info"))
	lg.info("scriptName")(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message info").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 0, logs.FilterMessage("message debug").FilterField(zap.String("scriptName", "scriptName")).Len())

	L = lua.NewState()
	L.Push(lua.LString("message debug"))
	lg.debug("scriptName")(L)
	assert.Equal(t, 1, logs.FilterMessage("message error").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message warn").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message info").FilterField(zap.String("scriptName", "scriptName")).Len())
	assert.Equal(t, 1, logs.FilterMessage("message debug").FilterField(zap.String("scriptName", "scriptName")).Len())
}

func TestNew(t *testing.T) {
	lgfunc := New(zap.NewNop())
	assert.IsType(t, &Log{}, lgfunc)
}

func TestName(t *testing.T) {
	l := &Log{}
	assert.Equal(t, "log", l.Name())
}

func TestStop(t *testing.T) {
	l := &Log{}
	assert.NoError(t, l.Stop())
}
