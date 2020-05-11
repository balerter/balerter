package loki

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	l, err := New(&config.DataSourceLoki{
		Name: "foo",
		URL:  "http://domain.com",
		BasicAuth: config.BasicAuth{
			Username: "user",
			Password: "secret",
		},
		Timeout: time.Second,
	}, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &Loki{}, l)
	assert.Equal(t, "loki.foo", l.name)
	assert.Equal(t, time.Second, l.timeout)
	assert.Equal(t, "http://domain.com", l.url.String())
}

func TestNewDefaultTimeout(t *testing.T) {
	l, err := New(&config.DataSourceLoki{}, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &Loki{}, l)
	assert.Equal(t, defaultTimeout, l.timeout)
}

func TestNewWrongURL(t *testing.T) {
	_, err := New(&config.DataSourceLoki{URL: "foobar\ncom"}, zap.NewNop())
	require.Error(t, err)
	require.Equal(t, "parse \"foobar\\ncom\": net/url: invalid control character in URL", err.Error())
}

func TestName(t *testing.T) {
	l := &Loki{name: "foo"}
	assert.Equal(t, "foo", l.Name())
}

func TestLoader(t *testing.T) {
	ch := &Loki{}

	f := ch.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}
