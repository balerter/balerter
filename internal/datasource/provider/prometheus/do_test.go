package prometheus

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
	m := &Prometheus{
		url: &url.URL{
			Host: "domain.com",
		},
	}

	u := m.sendRange("/foo", &queryRangeOptions{
		Start: "1",
		End:   "2",
		Step:  "3",
	})

	assert.Equal(t, "//domain.com/api/v1/query_range?end=2&query=%2Ffoo&start=1&step=3", u)
}

func Test_sendQuery(t *testing.T) {
	m := &Prometheus{
		url: &url.URL{
			Host: "domain.com",
		},
	}

	u := m.sendQuery("/foo", &queryQueryOptions{
		Time: "1",
	})

	assert.Equal(t, "//domain.com/api/v1/query?query=%2Ffoo&time=1", u)
}

type httpClientMock struct {
	mock.Mock
}

func (m *httpClientMock) CloseIdleConnections() {
	m.Called()
}

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
		Status:     "status1",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{"resultType":"vector","result":[]}}`))),
	}, nil)

	m := &Prometheus{
		logger:            zap.NewNop(),
		client:            hm,
		basicAuthUsername: "foo",
	}

	resp, err := m.send("domain.com/foo")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func Test_send_error(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(nil, fmt.Errorf("foo error"))

	m := &Prometheus{
		logger: zap.NewNop(),
		client: hm,
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "foo error", err.Error())
}

func Test_send_error_error_new_request(t *testing.T) {
	m := &Prometheus{
		logger: zap.NewNop(),
	}

	_, err := m.send("://://")
	require.Error(t, err)
	assert.Equal(t, "parse \"://://\": missing protocol scheme", err.Error())
}

func Test_send_bad_response_status(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(&http.Response{
		Status:     "status1",
		StatusCode: 0,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{"resultType":"vector","result":[]}}`))),
	}, nil)

	m := &Prometheus{
		logger:            zap.NewNop(),
		client:            hm,
		basicAuthUsername: "foo",
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "unexpected response code 0", err.Error())
}

func Test_send_bad_response_body(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(&http.Response{
		Status:     "status1",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(``))),
	}, nil)

	m := &Prometheus{
		logger:            zap.NewNop(),
		client:            hm,
		basicAuthUsername: "foo",
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}

func Test_send_bad_response_body2(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(&http.Response{
		Status:     "status1",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{}}`))),
	}, nil)

	m := &Prometheus{
		logger:            zap.NewNop(),
		client:            hm,
		basicAuthUsername: "foo",
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "unexpected value type \"\"", err.Error())
}

type badReader struct{}

func (br *badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("err1")
}

func Test_send_error_read_body(t *testing.T) {
	hm := &httpClientMock{}
	hm.On("Do", mock.Anything).Return(&http.Response{
		Status:     "status1",
		StatusCode: 200,
		Body:       ioutil.NopCloser(&badReader{}),
	}, nil)

	m := &Prometheus{
		logger:            zap.NewNop(),
		client:            hm,
		basicAuthUsername: "foo",
	}

	_, err := m.send("domain.com/foo")
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}
