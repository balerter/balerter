package alert

import (
	"bytes"
	"fmt"
	"time"
)

// Marshal returns JSON view if the Alert
func (a *Alert) Marshal() []byte {
	pattern := `{"name":"%s","level":"%s","level_num":%d,"count":%d,"last_change":"%s","start":"%s"}`

	return []byte(fmt.Sprintf(pattern, a.Name, a.Level.String(), a.Level, a.Count, a.LastChange.Format(time.RFC3339), a.Start.Format(time.RFC3339)))
}

// Marshal returns JSON view of the alerts slice
func (a Alerts) Marshal() []byte {
	bufs := make([][]byte, 0)

	for _, alert := range a {
		bufs = append(bufs, alert.Marshal())
	}

	buf := bytes.Join(bufs, []byte{','})

	return append([]byte{'['}, append(buf, ']')...)
}
