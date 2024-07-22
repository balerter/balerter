package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/balerter/balerter/internal/message"
)

type modelAlert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL"`
}

func newPromAlert() *modelAlert {
	return &modelAlert{
		Labels:       map[string]string{},
		Annotations:  map[string]string{},
		StartsAt:     time.Time{},
		EndsAt:       time.Time{},
		GeneratorURL: "",
	}
}

// Send a message to AlertManager
func (a *AlertManager) Send(mes *message.Message) error {
	promAlert := newPromAlert()

	// TODO (negasus): After refactoring with pass Alert to 'send' method, this condition should be refactoring
	if mes.Level == "success" {
		promAlert.EndsAt = time.Now()
	}

	promAlert.Labels["name"] = mes.AlertName
	promAlert.Annotations["description"] = mes.Text
	for k, v := range mes.Fields {
		promAlert.Labels[k] = v
	}

	data, err := json.Marshal([]*modelAlert{promAlert})
	if err != nil {
		return fmt.Errorf("error marshal prometheus alert, %w", err)
	}

	resp, err := a.whCore.Send(bytes.NewReader(data), mes)
	if err != nil {
		return fmt.Errorf("error send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	return nil
}
