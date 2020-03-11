package alerts

import (
	"encoding/json"
	"fmt"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"net/http"
	"time"
)

type resource struct {
	items []resourceItem
}

type resourceItem struct {
	Name      string `json:"name"`
	Level     string `json:"level"`
	Count     int    `json:"count"`
	UpdatedAt string `json:"updated_at"`
}

func newResource(data []*alertManager.AlertInfo) *resource {
	res := &resource{}

	for _, item := range data {
		i := resourceItem{
			Name:      item.Name,
			Level:     item.Level.String(),
			Count:     item.Count,
			UpdatedAt: item.LastChange.Format(time.RFC3339),
		}

		res.items = append(res.items, i)
	}

	return res
}

func (res *resource) render(rw http.ResponseWriter) error {
	data, err := json.Marshal(res.items)
	if err != nil {
		return fmt.Errorf("error marshal data, %w", err)
	}

	rw.Header().Add("Content-Type", "application/json")

	if _, err = fmt.Fprintf(rw, "%s", data); err != nil {
		return fmt.Errorf("error write response, %w", err)
	}

	return nil
}
