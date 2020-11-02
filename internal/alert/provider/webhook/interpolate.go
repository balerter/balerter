package webhook

import (
	"github.com/balerter/balerter/internal/alert/message"
	"strings"
)

const (
	MacrosLevel     = "$level"
	MacrosAlertName = "$alert_name"
	MacrosText      = "$text"
	MacrosImage     = "$image"
)

func interpolate(s string, m *message.Message) string {
	if m == nil {
		return s
	}

	return strings.NewReplacer(
		MacrosLevel, m.Level,
		MacrosAlertName, m.AlertName,
		MacrosText, m.Text,
		MacrosImage, m.Image,
	).Replace(s)
}
