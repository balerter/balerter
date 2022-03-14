package prometheus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	prometheusModels "github.com/balerter/balerter/internal/prometheus_models"
)

const (
	apiPrefix = "/api/v1"

	epQuery      = apiPrefix + "/query"
	epQueryRange = apiPrefix + "/query_range"
)

func (m *Prometheus) sendRange(query string, opts *queryRangeOptions) string {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	if opts.Start != "" {
		q.Add("start", opts.Start)
	}
	if opts.End != "" {
		q.Add("end", opts.End)
	}
	if opts.Step != "" {
		q.Add("step", opts.Step)
	}
	u.RawQuery = q.Encode()
	u.Path = epQueryRange

	return u.String()
}

func (m *Prometheus) sendQuery(query string, opts *queryQueryOptions) string {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	if opts.Time != "" {
		q.Add("time", opts.Time)
	}
	u.RawQuery = q.Encode()
	u.Path = epQuery

	return u.String()
}

func (m *Prometheus) send(u string) (prometheusModels.ModelValue, error) {
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiResp prometheusModels.APIResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data.Value, nil
}
