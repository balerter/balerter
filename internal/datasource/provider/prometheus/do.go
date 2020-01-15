package prometheus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/model"
	"io/ioutil"
	"net/http"
	"time"
)

func (m *Prometheus) do(query string) (model.Value, error) {

	u := fmt.Sprintf("%s/api/v1/query?query=%s", m.url, query)

	req, err := http.NewRequest(http.MethodGet, u, nil)
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
