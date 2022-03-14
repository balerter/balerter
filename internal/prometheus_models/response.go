package prometheus_models

import (
	"encoding/json"
	"fmt"
)

type APIResponse struct {
	Status    string      `json:"status"`
	Data      QueryResult `json:"data"`
	ErrorType string      `json:"errorType"`
	Error     string      `json:"error"`
	Warnings  []string    `json:"warnings,omitempty"`
}

type QueryResult struct {
	Type  string `json:"resultType"`
	Value ModelValue
}

func (qr *QueryResult) UnmarshalJSON(b []byte) error {
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
	case ValVector:
		var vv Vector
		err = json.Unmarshal(v.Result, &vv)
		qr.Value = vv

	case ValMatrix:
		var mv Matrix
		err = json.Unmarshal(v.Result, &mv)
		qr.Value = mv

	case ValStreams:
		var mv Streams
		err = json.Unmarshal(v.Result, &mv)
		qr.Value = mv

	default:
		err = fmt.Errorf("unexpected value type %q", v.Type)
	}
	return err
}
