package alertmanagerreceiver

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type webHookCoreMock struct {
	mock.Mock
}

func (mm *webHookCoreMock) Send(body io.Reader, m *message.Message) (*http.Response, error) {
	args := mm.Called(body, m)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSend_error(t *testing.T) {
	m := &webHookCoreMock{}

	am := &AMReceiver{
		whCore: m,
	}

	m.On("Send", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("err1"))

	mes := &message.Message{
		Level: "success",
	}

	err := am.Send(mes)
	require.Error(t, err)
	assert.Equal(t, "error send request, err1", err.Error())
}

func TestSend_error_status_code(t *testing.T) {
	m := &webHookCoreMock{}

	am := &AMReceiver{
		whCore: m,
	}

	resp := &http.Response{StatusCode: 0, Body: ioutil.NopCloser(bytes.NewBuffer(nil))}

	m.On("Send", mock.Anything, mock.Anything).Return(resp, nil)

	mes := &message.Message{
		Level: "error",
	}

	err := am.Send(mes)
	require.Error(t, err)
	assert.Equal(t, "unexpected response status code 0", err.Error())
}

func TestSend(t *testing.T) {
	m := &webHookCoreMock{}

	am := &AMReceiver{
		whCore: m,
	}

	resp := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(nil))}

	m.On("Send", mock.Anything, mock.Anything).Return(resp, nil)

	mes := &message.Message{
		Level: "success",
	}

	err := am.Send(mes)
	require.NoError(t, err)
}
