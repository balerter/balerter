package prometheus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/model"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (m *Prometheus) query(L *lua.LState) int {
	query := strings.TrimSpace(L.Get(1).String())
	if query == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("query must be not empty"))
		return 2
	}

	m.logger.Debug("call prometheus query", zap.String("name", m.name), zap.String("query", query))

	v, err := m.do(query)
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		L.Push(lua.LNil)
		L.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	fmt.Printf("\n\n\n%+v\n\n\n", v)

	L.Push(lua.LNil)
	L.Push(lua.LNil)

	return 2
}

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
