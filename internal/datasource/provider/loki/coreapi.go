package loki

import (
	"fmt"
	"net/http"
	"strconv"
)

func (m *Loki) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	var u string

	switch method {
	case "query":
		var opts queryOptions
		if v, ok := params["time"]; ok {
			opts.Time = v
		}
		if v, ok := params["limit"]; ok {
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid limit value: %s", v)
			}
			opts.Limit = n
		}
		if v, ok := params["direction"]; ok {
			if err := directionValidate(v); err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid direction value: %s", v)
			}
			opts.Direction = v
		}

		u = m.sendQuery(string(body), &opts)
	case "range":
		var opts rangeOptions
		if v, ok := params["limit"]; ok {
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid limit value: %s", v)
			}
			opts.Limit = n
		}
		if v, ok := params["start"]; ok {
			opts.Start = v
		}
		if v, ok := params["end"]; ok {
			opts.End = v
		}
		if v, ok := params["step"]; ok {
			opts.Step = v
		}
		if v, ok := params["direction"]; ok {
			if err := directionValidate(v); err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("invalid direction value: %s", v)
			}
			opts.Direction = v
		}
		u = m.sendRange(string(body), &opts)
	default:
		return nil, http.StatusNotFound, fmt.Errorf("unknown method: %q", method)
	}

	resp, errSend := m.send(u)
	if errSend != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error send query to loki: %w", errSend)
	}

	return resp, 0, nil
}
