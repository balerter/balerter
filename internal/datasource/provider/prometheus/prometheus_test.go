package prometheus

import (
	"github.com/balerter/balerter/internal/config/common"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	m, err := New(prometheus.Prometheus{
		Name:      "prom1",
		URL:       "http://domain.com",
		BasicAuth: common.BasicAuth{},
		Timeout:   0,
	}, zap.NewNop())

	require.NoError(t, err)

	assert.Equal(t, "prometheus.prom1", m.name)
}

func TestNew_fail_url(t *testing.T) {
	_, err := New(prometheus.Prometheus{
		Name:      "prom1",
		URL:       "$% a.a",
		BasicAuth: common.BasicAuth{},
		Timeout:   0,
	}, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "parse \"$% a.a\": invalid URL escape \"% a\"", err.Error())
}

func TestName(t *testing.T) {
	m := &Prometheus{name: "prom1"}
	assert.Equal(t, "prom1", m.Name())
}

func TestStop(t *testing.T) {
	mm := &httpClientMock{}
	mm.On("CloseIdleConnections").Return()

	m := &Prometheus{
		client: mm,
	}

	err := m.Stop()
	require.NoError(t, err)

	mm.AssertCalled(t, "CloseIdleConnections")
	mm.AssertExpectations(t)
}

func TestModuleName(t *testing.T) {
	assert.Equal(t, "prometheus.foo", ModuleName("foo"))
}

func TestMethods(t *testing.T) {
	a := Methods()
	assert.Equal(t, 2, len(a))
	assert.Equal(t, "query", a[0])
	assert.Equal(t, "range", a[1])
}

func TestPrometheus_GetLoader(t *testing.T) {
	p := &Prometheus{}
	loader := p.GetLoader(nil)

	luaState := lua.NewState()

	n := loader(luaState)
	assert.Equal(t, 1, n)

	v := luaState.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}
