package models

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/loki/pkg/logql/stats"
)

const (
	ResultTypeStream = "streams"
	ResultTypeScalar = "scalar"
	ResultTypeVector = "vector"
	ResultTypeMatrix = "matrix"
)

type ResultType string

type ResultValue interface {
	Type() ResultType
}

type QueryResponseData struct {
	ResultType ResultType  `json:"resultType"`
	Result     ResultValue `json:"result"`
	//Statistics stats.Result `json:"stats"`
}

type QueryResponse struct {
	Status string            `json:"status"`
	Data   QueryResponseData `json:"data"`
}

func (q *QueryResponseData) UnmarshalJSON(data []byte) error {
	unmarshal := struct {
		Type       ResultType      `json:"resultType"`
		Result     json.RawMessage `json:"result"`
		Statistics stats.Result    `json:"stats"`
	}{}

	err := json.Unmarshal(data, &unmarshal)
	if err != nil {
		return err
	}

	var value ResultValue

	// unmarshal results
	switch unmarshal.Type {
	case ResultTypeStream:
		var s Streams
		err = json.Unmarshal(unmarshal.Result, &s)
		value = s
	//case ResultTypeMatrix:
	//	var m Matrix
	//	err = json.Unmarshal(unmarshal.Result, &m)
	//	value = m
	//case ResultTypeVector:
	//	var v Vector
	//	err = json.Unmarshal(unmarshal.Result, &v)
	//	value = v
	//case ResultTypeScalar:
	//	var v Scalar
	//	err = json.Unmarshal(unmarshal.Result, &v)
	//	value = v
	default:
		return fmt.Errorf("unknown type: %s", unmarshal.Type)
	}

	if err != nil {
		return err
	}

	q.ResultType = unmarshal.Type
	q.Result = value
	//q.Statistics = unmarshal.Statistics

	return nil
}
