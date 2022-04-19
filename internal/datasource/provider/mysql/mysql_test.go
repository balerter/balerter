package mysql

import (
	"fmt"
	"testing"

	"github.com/balerter/balerter/internal/config/datasources/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func TestNew_ErrorConnect(t *testing.T) {
	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return nil, fmt.Errorf("err1")
	}

	cfg := mysql.Mysql{}

	_, err := New(cfg, mockConnFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestName(t *testing.T) {
	p := &MySQL{name: "Foo"}
	assert.Equal(t, "Foo", p.Name())
}

func TestGetLoader(t *testing.T) {
	p := &MySQL{}

	loader := p.GetLoader(nil)

	luaState := lua.NewState()

	n := loader(luaState)
	assert.Equal(t, 1, n)

	v := luaState.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestModuleName(t *testing.T) {
	assert.Equal(t, "mysql.Foo", ModuleName("Foo"))
}
