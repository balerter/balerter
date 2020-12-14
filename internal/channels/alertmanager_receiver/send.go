package alertmanagerreceiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/balerter/balerter/internal/message"
	amWebhook "github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/common/model"
	"net/http"
	"time"
)

func (a *AMReceiver) Send(mes *message.Message) error {
	data := &template.Data{
		Receiver:          "balerter",
		Status:            string(model.AlertResolved),
		Alerts:            nil,
		GroupLabels:       nil,
		CommonLabels:      nil,
		CommonAnnotations: template.KV{"name": mes.AlertName, "description": mes.Text},
		ExternalURL:       "",
	}

	alrt := template.Alert{
		Status:       string(model.AlertResolved),
		Labels:       template.KV{"name": mes.AlertName},
		Annotations:  template.KV{"name": mes.AlertName, "description": mes.Text},
		StartsAt:     time.Time{},
		EndsAt:       time.Now(),
		GeneratorURL: "",
		Fingerprint:  "",
	}

	// TODO (negasus): After refactoring with pass Alert to 'send' method, this condition should be refactoring
	if mes.Level == "error" {
		data.Status = string(model.AlertFiring)
		alrt.Status = string(model.AlertFiring)
		alrt.StartsAt = time.Now()
		alrt.EndsAt = time.Time{}
	}

	data.Alerts = append(data.Alerts, alrt)

	amMes := &amWebhook.Message{
		Version:  "4",
		GroupKey: "",
		Data:     data,
	}

	buf, err := json.Marshal(amMes)
	if err != nil {
		return fmt.Errorf("error marshal message, %w", err)
	}

	resp, err := a.whCore.Send(bytes.NewReader(buf), mes)
	if err != nil {
		return fmt.Errorf("error send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	return nil
}
