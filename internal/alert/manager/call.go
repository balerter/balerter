package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/metrics"
	"go.uber.org/zap"
)

func (m *Manager) Call(name string, level alert.Level, text string, options *alert.Options) error {
	m.logger.Debug("alertManager call",
		zap.String("name", name),
		zap.Int("level", int(level)),
		zap.String("text", text),
		zap.Any("options", options),
	)

	metrics.SetAlertLevel(name, level)

	a, err := m.engine.Alert().GetOrNew(name)
	if err != nil {
		return fmt.Errorf("error get alert %s, %w", name, err)
	}

	if a.HasLevel(level) {
		a.Inc()

		if !options.Quiet && options.Repeat > 0 && a.Count()%options.Repeat == 0 {
			m.Send(level.String(), name, text, options.Channels, options.Fields, options.Image)
		}

		return nil
	}

	a.UpdateLevel(level)

	if !options.Quiet {
		m.Send(level.String(), name, text, options.Channels, options.Fields, options.Image)
	}

	if err := m.engine.Alert().Release(a); err != nil {
		m.logger.Error("error release alert", zap.Error(err))
	}

	return nil
}
