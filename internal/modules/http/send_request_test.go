package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHTTP_sendRequest_error_new_request(t *testing.T) {
	h := &HTTP{}
	args := &requestArgs{
		URI: "://://",
	}
	_, err := h.sendRequest(args)
	require.Error(t, err)
	assert.Equal(t, "error build request, parse \"://://\": missing protocol scheme", err.Error())
}

func TestHTTP_sendRequest_error_do_request(t *testing.T) {
	hm := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	h := &HTTP{
		createClientFunc: func(_ time.Duration, _ bool) httpClient { return hm },
	}
	args := &requestArgs{}
	_, err := h.sendRequest(args)
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

type badReader struct{}

func (br *badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("err1")
}

func TestHTTP_sendRequest_error_read_response_body(t *testing.T) {
	hm := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{}
			r.Body = io.NopCloser(&badReader{})
			return r, nil
		},
	}

	h := &HTTP{
		createClientFunc: func(_ time.Duration, _ bool) httpClient { return hm },
	}
	args := &requestArgs{}
	_, err := h.sendRequest(args)
	require.Error(t, err)
	assert.Equal(t, "error read body, err1", err.Error())
}

func TestHTTP_sendRequest_multiple_response_header_values(t *testing.T) {
	hm := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{
				StatusCode: 10,
				Header:     map[string][]string{"foo": {"bar", "baz"}},
			}
			r.Body = io.NopCloser(bytes.NewBuffer([]byte("foo")))
			return r, nil
		},
	}

	core, logs := observer.New(zap.DebugLevel)

	h := &HTTP{
		createClientFunc: func(_ time.Duration, _ bool) httpClient { return hm },
		logger:           zap.New(core),
	}
	args := &requestArgs{
		Headers: map[string]string{"a": "b"},
	}
	resp, err := h.sendRequest(args)
	require.NoError(t, err)

	assert.Equal(t, 1, logs.FilterMessage("the response header has multiple values").Len())
	assert.Equal(t, 10, resp.StatusCode)
	assert.Equal(t, "foo", string(resp.Body))
	assert.Equal(t, 1, len(resp.Headers))
	v, ok := resp.Headers["foo"]
	assert.True(t, ok)
	assert.Equal(t, "baz", v)
}
