package twiliovoice

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

func TestTwilioVoice_Send_error_new_request(t *testing.T) {
	tw := &TwilioVoice{
		apiPrefix: "://://",
	}

	err := tw.Send(&message.Message{})
	require.Error(t, err)
	assert.Equal(t, "parse \"://:///Accounts//Calls.json\": missing protocol scheme", err.Error())
}

func TestTwilioVoice_Send_error_send_request(t *testing.T) {
	c := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	tw := &TwilioVoice{
		client: c,
	}

	err := tw.Send(&message.Message{})
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

type badReader struct{}

func (br *badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("err1")
}

func TestTwilioVoice_Send_error_read_response_body(t *testing.T) {
	c := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body: io.NopCloser(&badReader{}),
			}
			return r, nil
		},
	}

	tw := &TwilioVoice{
		client: c,
	}

	err := tw.Send(&message.Message{})
	require.Error(t, err)
	assert.Equal(t, "error read response body, err1", err.Error())
}

func TestTwilioVoice_Send_bad_response_code(t *testing.T) {
	c := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
				StatusCode: http.StatusOK,
			}
			return r, nil
		},
	}

	tw := &TwilioVoice{
		client: c,
		logger: zap.NewNop(),
	}

	err := tw.Send(&message.Message{})
	require.Error(t, err)
	assert.Equal(t, "unexpected status code 200", err.Error())
}

func TestTwilioVoice_Send(t *testing.T) {
	c := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
				StatusCode: http.StatusCreated,
			}
			return r, nil
		},
	}

	tw := &TwilioVoice{
		client: c,
		logger: zap.NewNop(),
	}

	err := tw.Send(&message.Message{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(c.DoCalls()))
}
