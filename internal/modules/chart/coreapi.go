package chart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type coreapiRequest struct {
	Title  string       `json:"title"`
	Series []DataSeries `json:"series"`
}

func (ch *Chart) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method != "render" {
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
	}
	var r = coreapiRequest{}
	errUnmarshal := json.Unmarshal(body, &r)
	if errUnmarshal != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unmarshal error, %w", errUnmarshal)
	}

	buf := bytes.NewBuffer(nil)

	errRender := ch.Render(r.Title, &Data{Title: r.Title, Series: r.Series}, buf)
	if errRender != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("render error, %w", errRender)
	}

	return buf.Bytes(), 0, nil
}
