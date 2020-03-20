package kv

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type resource struct {
	items []resourceItem
}

type resourceItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func newResource(data map[string]string) *resource {
	res := &resource{
		items: []resourceItem{},
	}

	for name, value := range data {
		i := resourceItem{
			Name:  name,
			Value: value,
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
