package metrics

import (
	"fmt"

	"github.com/balerter/balerter/internal/alert"

	"github.com/VictoriaMetrics/metrics"
)

// SetAlertLevel sets alert level
func SetAlertLevel(name string, level alert.Level) {
	metricsName := fmt.Sprintf("balerter_alert_status{name=%q}", name)
	metrics.GetOrCreateGauge(metricsName, func() float64 {
		return float64(level)
	})
}

// GetAlertLevel returns alert level from the metric
// use for testing purposes
func GetAlertLevel(name string) (float64, error) {
	metricsName := fmt.Sprintf("balerter_alert_status{name=%q}", name)

	v := metrics.GetOrCreateGauge(metricsName, func() float64 {
		return -1
	}).Get()

	return v, nil
}

// SetVersion updates data for metrics metricInfoVersion
func SetVersion(version string) {
	name := fmt.Sprintf("balerter_info_version{version=%q}", version)
	metrics.GetOrCreateGauge(name, func() float64 {
		return 1
	})
}

// SetScriptsActive updates data for metrics metricScripts
func SetScriptsActive(name string, active bool) {
	metricsName := fmt.Sprintf("balerter_scripts_active{name=%q}", name)
	if active {
		metrics.GetOrCreateGauge(metricsName, func() float64 {
			return 1
		})
		return
	}
	metrics.GetOrCreateGauge(metricsName, func() float64 {
		return 0
	})
}
