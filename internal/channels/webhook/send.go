package webhook

import (
	"github.com/balerter/balerter/internal/message"
	"strings"
)

func (w *Webhook) Send(m *message.Message) error {
	body := interpolate(w.body, m)

	resp, err := w.whCore.Send(strings.NewReader(body), m)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// todo: maybe check response code (add to config)

	return nil
}
