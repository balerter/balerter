package loki

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func Test_sendRange(t *testing.T) {
	m := &Loki{
		url: &url.URL{
			Host: "domain.com",
		},
	}

	u := m.sendRange("", &rangeOptions{
		Limit:     10,
		Start:     "1",
		End:       "2",
		Step:      "3",
		Direction: "4",
	})

	assert.Equal(t, "//domain.com/loki/api/v1/query_range?direction=4&end=2&limit=10&query=&start=1&step=3", u)
}

func Test_sendQuery(t *testing.T) {
	m := &Loki{
		url: &url.URL{
			Host: "domain.com",
		},
	}

	u := m.sendQuery("", &queryOptions{
		Time:      "1",
		Limit:     10,
		Direction: "2",
	})

	assert.Equal(t, "//domain.com/loki/api/v1/query?direction=2&limit=10&query=&time=1", u)
}

type httpClientMock struct {
	mock.Mock
}

func (m *httpClientMock) CloseIdleConnections() {}
func (m *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return r.(*http.Response), args.Error(1)
}

func Test_send(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(&http.Response{
		Status: "status1",
		Body:   ioutil.NopCloser(bytes.NewBuffer([]byte(`{}`))),
	}, nil)

	m := &Loki{
		logger: zap.NewNop(),
		client: hm,
	}

	resp, err := m.send("domain.com/foo")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func Test_send_error(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(nil, fmt.Errorf("foo error"))

	m := &Loki{
		logger: zap.NewNop(),
		client: hm,
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "foo error", err.Error())
}
