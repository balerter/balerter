package clickhouse

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestName(t *testing.T) {
	ch := &Clickhouse{name: "foo"}
	assert.Equal(t, "foo", ch.Name())
}

func TestLoader(t *testing.T) {
	ch := &Clickhouse{}

	f := ch.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("query")))
}
