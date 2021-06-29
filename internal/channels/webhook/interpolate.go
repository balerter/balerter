package webhook

import (
	"github.com/balerter/balerter/internal/message"
	"strings"
)

const (
	macrosLevel     = "$level"
	macrosAlertName = "$alert_name"
	macrosText      = "$text"
	macrosImage     = "$image"
)

func interpolate(s string, m *message.Message) string {
	if m == nil {
		return s
	}

	return strings.NewReplacer(
		macrosLevel, m.Level,
		macrosAlertName, m.AlertName,
		macrosText, m.Text,
		macrosImage, m.Image,
	).Replace(s)
}
