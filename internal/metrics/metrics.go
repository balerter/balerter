package metrics

import (
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

	metricScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "balerter",
		Subsystem:   "scripts",
		Name:        "active",
		Help:        "Balerter active scripts list",
		ConstLabels: nil,
	}, []string{"name"})
)

func Register(logger *zap.Logger) {
	logger.Debug("register metrics")

	if err := prometheus.Register(metricInfoVersion); err != nil {
		logger.Error("error register metrics", zap.String("name", "infoVersion"), zap.Error(err))
	}
	if err := prometheus.Register(metricScripts); err != nil {
		logger.Error("error register metrics", zap.String("name", "scripts"), zap.Error(err))
	}
}

func SetVersion(version string) {
	metricInfoVersion.WithLabelValues(version).Inc()
}

func SetScriptsActive(name string, active bool) {
	if active {
		metricScripts.WithLabelValues(name).Set(1)
		return
	}
	metricScripts.WithLabelValues(name).Set(0)
}
