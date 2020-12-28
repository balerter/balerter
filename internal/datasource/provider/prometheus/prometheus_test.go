package prometheus

import (
	"github.com/balerter/balerter/internal/config/common"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	m, err := New(&prometheus.Prometheus{
		Name:      "prom1",
		URL:       "http://domain.com",
		BasicAuth: common.BasicAuth{},
		Timeout:   0,
	}, zap.NewNop())

	require.NoError(t, err)

	assert.Equal(t, "prometheus.prom1", m.name)
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
