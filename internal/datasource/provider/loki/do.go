package loki

import (
	"context"
	"encoding/base64"
	"encoding/json"
	lokihttp "github.com/grafana/loki/pkg/loghttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPrefix = "/loki/api/v1"

	epQuery      = apiPrefix + "/query"
	epQueryRange = apiPrefix + "/query_range"
)

func (m *Loki) sendRange(query string, opts *rangeOptions) (*lokihttp.QueryResponse, error) {
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

	return m.send(&u)
}

func (m *Loki) sendQuery(query string, opts *queryOptions) (*lokihttp.QueryResponse, error) {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	if opts.Time != "" {
		q.Add("time", opts.Time)
	}
	u.RawQuery = q.Encode()
	u.Path = epQuery

	return m.send(&u)
}

func (m *Loki) send(u *url.URL) (*lokihttp.QueryResponse, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if m.basicAuthUsername != "" {
		ba := base64.StdEncoding.EncodeToString([]byte(m.basicAuthUsername + ":" + m.basicAuthPassword))
		req.Header.Add("Authorization", "Basic "+ba)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apires *lokihttp.QueryResponse

	err = json.Unmarshal(body, &apires)
	if err != nil {
		return nil, err
	}

	return apires, nil
}
