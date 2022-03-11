package loki

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/balerter/balerter/internal/datasource/provider/loki/models"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	apiPrefix = "/loki/api/v1"

	epQuery      = apiPrefix + "/query"
	epQueryRange = apiPrefix + "/query_range"
)

func (m *Loki) sendRange(query string, opts *rangeOptions) string {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	if opts.Limit > 0 {
		q.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Start != "" {
		q.Add("start", opts.Start)
	}
	if opts.End != "" {
		q.Add("end", opts.End)
	}
	if opts.Step != "" {
		q.Add("step", opts.Step)
	}
	if opts.Direction != "" {
		q.Add("direction", opts.Direction)
	}
	u.RawQuery = q.Encode()
	u.Path = epQueryRange

	return u.String()
}

func (m *Loki) sendQuery(query string, opts *queryOptions) string {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	if opts.Limit > 0 {
		q.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Time != "" {
		q.Add("time", opts.Time)
	}
	if opts.Direction != "" {
		q.Add("direction", opts.Direction)
	}
	u.RawQuery = q.Encode()
	u.Path = epQuery

	return u.String()
}

func (m *Loki) send(u string) (*models.QueryResponse, error) {
	m.logger.Debug("request to loki", zap.String("url", u))

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	if m.basicAuthUsername != "" {
		ba := base64.StdEncoding.EncodeToString([]byte(m.basicAuthUsername + ":" + m.basicAuthPassword))
		req.Header.Add("Authorization", "Basic "+ba)
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiResp *models.QueryResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}
