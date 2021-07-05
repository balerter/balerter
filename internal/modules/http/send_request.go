package http

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func (h *HTTP) sendRequest(args *requestArgs) (*response, error) {
	req, err := http.NewRequest(args.Method, args.URI, bytes.NewReader(args.Body))
	if err != nil {
		return nil, fmt.Errorf("error build request, %w", err)
	}

	for name, value := range args.Headers {
		req.Header.Set(name, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := newResponse()

	res.StatusCode = resp.StatusCode

	res.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read body, %w", err)
	}

	for name, values := range resp.Header {
		if len(values) > 1 {
			h.logger.Debug("the response header has multiple values", zap.String("name", name), zap.Strings("values", values))
		}
		for _, value := range values {
			res.Headers[name] = value
		}
	}

	return res, nil
}
