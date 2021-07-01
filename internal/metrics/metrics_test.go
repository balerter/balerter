package metrics

import (
	"github.com/balerter/balerter/internal/alert"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestSetAlertLevel(t *testing.T) {
	SetAlertLevel("foo", alert.LevelWarn)

	v, err := GetAlertLevel("foo")
	require.NoError(t, err)
	assert.Equal(t, float64(alert.LevelWarn), v)
}

func TestSetVersion(t *testing.T) {
	SetVersion("1")

	g, err := metricInfoVersion.GetMetricWithLabelValues("1")
	require.NoError(t, err)

	d := &dto.Metric{}
	err = g.Write(d)
	require.NoError(t, err)

	assert.Equal(t, float64(1), *d.Gauge.Value)
}

func TestSetScriptsActive_true(t *testing.T) {
	SetScriptsActive("foo", true)

	g, err := metricScripts.GetMetricWithLabelValues("foo")
	require.NoError(t, err)

	d := &dto.Metric{}
	err = g.Write(d)
	require.NoError(t, err)

	assert.Equal(t, float64(1), *d.Gauge.Value)
}

func TestSetScriptsActive_false(t *testing.T) {
	SetScriptsActive("foo", false)

	g, err := metricScripts.GetMetricWithLabelValues("foo")
	require.NoError(t, err)

	d := &dto.Metric{}
	err = g.Write(d)
	require.NoError(t, err)

	assert.Equal(t, float64(0), *d.Gauge.Value)
}

func TestRegister(t *testing.T) {
	Register(zap.NewNop())
}
