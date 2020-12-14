package metrics

import (
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/prometheus/client_golang/prometheus"
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

	metricAlert = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "alert",
		Name:        "status",
		Help:        "Balerter alerts status",
		ConstLabels: nil,
	}, []string{"name"})

	metricScriptsCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "scripts",
		Name:        "count",
		Help:        "Balerter scripts count",
		ConstLabels: nil,
	})
)

func Register(logger *zap.Logger) {
	logger.Debug("register metrics")

	if err := prometheus.Register(metricInfoVersion); err != nil {
		logger.Error("error register metrics", zap.String("name", "infoVersion"), zap.Error(err))
	}
	if err := prometheus.Register(metricAlert); err != nil {
		logger.Error("error register metrics", zap.String("name", "alert"), zap.Error(err))
	}
	if err := prometheus.Register(metricScriptsCount); err != nil {
		logger.Error("error register metrics", zap.String("name", "scriptsCount"), zap.Error(err))
	}
}

func SetVersion(version string) {
	metricInfoVersion.WithLabelValues(version).Inc()
}

func SetAlertLevel(alertName string, level alert2.Level) {
	metricAlert.WithLabelValues(alertName).Set(float64(level))
}

func SetScriptsCount(count int) {
	metricScriptsCount.Set(float64(count))
}
