package prometheus

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestPrometheus_doQuery_empty_query(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	n := m.doQuery(luaState)
	assert.Equal(t, 2, n)
}
