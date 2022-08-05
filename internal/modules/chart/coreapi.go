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

func (ch *Chart) CoreApiHandler(req []string, body []byte) (any, int, error) {
	if len(req) != 1 {
		return nil, http.StatusBadRequest, fmt.Errorf("wrong request length")
	}
	if req[0] != "render" {
		return nil, http.StatusBadRequest, fmt.Errorf("wrong request, unknown method: " + req[0])
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
