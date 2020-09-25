package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/prometheus/common/model"
	"net/http"
	"time"
)

func newPromAlert() *model.Alert {
	return &model.Alert{
		Labels:       map[model.LabelName]model.LabelValue{},
		Annotations:  map[model.LabelName]model.LabelValue{},
		StartsAt:     time.Time{},
		EndsAt:       time.Time{},
		GeneratorURL: "",
	}
}

func (a *AlertManager) Send(mes *message.Message) error {
	promAlert := newPromAlert()

	// TODO (negasus): After refactoring with pass Alert to 'send' method, this condition should be refactoring
	if mes.Level == "success" {
		promAlert.EndsAt = time.Now()
	}

	promAlert.Labels["name"] = model.LabelValue(mes.AlertName)
	promAlert.Annotations["description"] = model.LabelValue(mes.Text)

	data, err := json.Marshal([]*model.Alert{promAlert})
	if err != nil {
		return fmt.Errorf("error marshal prometheus alert, %w", err)
	}

	resp, err := http.Post(a.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error send request to %s, %w", a.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	return nil
}
