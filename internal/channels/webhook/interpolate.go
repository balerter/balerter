package webhook

import (
	"fmt"
	"github.com/balerter/balerter/internal/message"
	"strings"
)

const (
	macrosLevel     = "$level"
	macrosAlertName = "$alert_name"
	macrosText      = "$text"
	macrosImage     = "$image"
	macrosFields    = "$fields"
)

func interpolate(s string, m *message.Message) string {
	if m == nil {
		return s
	}

	var fields string
	if len(m.Fields) > 0 {
		var f []string
		for k, v := range m.Fields {
			f = append(f, fmt.Sprintf("%s=%s", k, v))
		}
		fields = strings.Join(f, ",")
	}

	return strings.NewReplacer(
		macrosLevel, m.Level,
		macrosAlertName, m.AlertName,
		macrosText, m.Text,
		macrosImage, m.Image,
		macrosFields, fields,
	).Replace(s)
}
