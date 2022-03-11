package prometheus

import (
	"encoding/json"
	"fmt"

	"github.com/balerter/balerter/internal/datasource/provider/prometheus/models"
)

type apiResponse struct {
	Status    string      `json:"status"`
	Data      queryResult `json:"data"`
	ErrorType string      `json:"errorType"`
	Error     string      `json:"error"`
	Warnings  []string    `json:"warnings,omitempty"`
}

type queryResult struct {
	Type string `json:"resultType"`
	v    models.ModelValue
}

func (qr *queryResult) UnmarshalJSON(b []byte) error {
	v := struct {
		Type   string          `json:"resultType"`
		Result json.RawMessage `json:"result"`
	}{}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	qr.Type = v.Type

	switch v.Type {
	case models.ValVector:
		var vv models.Vector
		err = json.Unmarshal(v.Result, &vv)
		qr.v = vv

	case models.ValMatrix:
		var mv models.Matrix
		err = json.Unmarshal(v.Result, &mv)
		qr.v = mv

	default:
		err = fmt.Errorf("unexpected value type %q", v.Type)
	}
	return err
}
