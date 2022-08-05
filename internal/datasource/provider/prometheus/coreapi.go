package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (m *Prometheus) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	var u string

	switch method {
	case "query":
		opts := &queryQueryOptions{}
		if v, ok := params["time"]; ok {
			opts.Time = v
		}
		u = m.sendQuery(string(body), opts)
	case "range":
		opts := &queryRangeOptions{}
		if v, ok := params["start"]; ok {
			opts.Start = v
		}
		if v, ok := params["end"]; ok {
			opts.End = v
		}
		if v, ok := params["step"]; ok {
			opts.Step = v
		}
		u = m.sendRange(string(body), opts)
	default:
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
	}

	res, errDo := m.send(u)
	if errDo != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error sending request: %v", errDo)
	}

	resp, errMarshal := json.Marshal(res)
	if errMarshal != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error marshaling response: %v", errMarshal)
	}

	return resp, 0, nil
}
