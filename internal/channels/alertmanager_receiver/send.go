package alertmanagerreceiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/balerter/balerter/internal/message"
)

type AlertStatus string

const (
	AlertFiring   AlertStatus = "firing"
	AlertResolved AlertStatus = "resolved"
)

type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

type KV map[string]string

type Alerts []Alert

type Data struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

type Message struct {
	*Data

	// The protocol version.
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts uint64 `json:"truncatedAlerts"`
}

// Send message to the channel
func (a *AMReceiver) Send(mes *message.Message) error {
	data := &Data{
		Receiver:          "balerter",
		Status:            string(AlertResolved),
		Alerts:            nil,
		GroupLabels:       nil,
		CommonLabels:      nil,
		CommonAnnotations: KV{"name": mes.AlertName, "description": mes.Text},
		ExternalURL:       "",
	}

	alrt := Alert{
		Status:       string(AlertResolved),
		Labels:       KV{"name": mes.AlertName},
		Annotations:  KV{"name": mes.AlertName, "description": mes.Text},
		StartsAt:     time.Time{},
		EndsAt:       time.Now(),
		GeneratorURL: "",
		Fingerprint:  "",
	}

	// TODO (negasus): After refactoring with pass Alert to 'send' method, this condition should be refactoring
	if mes.Level == "error" {
		data.Status = string(AlertFiring)
		alrt.Status = string(AlertFiring)
		alrt.StartsAt = time.Now()
		alrt.EndsAt = time.Time{}
	}

	data.Alerts = append(data.Alerts, alrt)

	amMes := &Message{
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
