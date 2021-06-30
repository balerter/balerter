package metrics

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"go.uber.org/zap"
)

var (
	metricInfoVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "info",
		Name:        "version",
		Help:        "Information about the Balerter environment",
		ConstLabels: nil,
	}, []string{"version"})

	metricScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "scripts",
		Name:        "active",
		Help:        "Balerter active scripts list",
		ConstLabels: nil,
	}, []string{"name"})

	metricAlert = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "alert",
		Name:        "status",
		Help:        "Balerter alerts status",
		ConstLabels: nil,
	}, []string{"name"})
)

// Register metrics for expose
func Register(logger *zap.Logger) {
	logger.Debug("register metrics")

	if err := prometheus.Register(metricInfoVersion); err != nil {
		logger.Error("error register metrics", zap.String("name", "infoVersion"), zap.Error(err))
	}
	if err := prometheus.Register(metricScripts); err != nil {
		logger.Error("error register metrics", zap.String("name", "scripts"), zap.Error(err))
	}
	if err := prometheus.Register(metricAlert); err != nil {
		logger.Error("error register metrics", zap.String("name", "alert"), zap.Error(err))
	}
}

// SetAlertLevel sets alert level
func SetAlertLevel(name string, level alert.Level) {
	metricAlert.WithLabelValues(name).Set(float64(level))
}

// GetAlertLevel returns alert level from the metric
// use for testing purposes
func GetAlertLevel(name string) (float64, error) {
	g, err := metricAlert.GetMetricWithLabelValues(name)
	if err != nil {
		return 0, err
	}

	d := &dto.Metric{}
	err = g.Write(d)
	if err != nil {
		return 0, err
	}

	return *d.Gauge.Value, nil
}

// SetVersion updates data for metrics metricInfoVersion
func SetVersion(version string) {
	metricInfoVersion.WithLabelValues(version).Inc()
}

// SetScriptsActive updates data for metrics metricScripts
func SetScriptsActive(name string, active bool) {
	if active {
		metricScripts.WithLabelValues(name).Set(1)
		return
	}
	metricScripts.WithLabelValues(name).Set(0)
}
