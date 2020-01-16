package prometheus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/prometheus/common/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	//statusAPIError = 422

	apiPrefix = "/api/v1"

	epQuery      = apiPrefix + "/query"
	epQueryRange = apiPrefix + "/query_range"
)

func (m *Prometheus) sendRange(query string, opts queryRangeOptions) (model.Value, error) {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	q.Add("start", opts.Start)
	q.Add("end", opts.End)
	q.Add("step", opts.Step)
	u.RawQuery = q.Encode()
	u.Path = epQueryRange

	return m.send(&u)
}

func (m *Prometheus) sendQuery(query string) (model.Value, error) {
	u := *m.url

	q := &url.Values{}
	q.Add("query", query)
	u.RawQuery = q.Encode()
	u.Path = epQuery

	return m.send(&u)
}

func (m *Prometheus) send(u *url.URL) (model.Value, error) {
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

	var apires apiResponse
	var qres queryResult

	err = json.Unmarshal(body, &apires)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(apires.Data, &qres)
	if err != nil {
		return nil, err
	}

	return qres.v, nil
}
