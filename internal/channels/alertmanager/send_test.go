package alertmanager

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/balerter/balerter/internal/message"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func Test_newPromAlert(t *testing.T) {
	a := newPromAlert()
	assert.IsType(t, &modelAlert{}, a)
}

func TestSend_error_send(t *testing.T) {
	m := &webHookCoreMock{}

	a := &AlertManager{
		whCore: m,
		logger: zap.NewNop(),
	}

	m.On("Send", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("err1"))

	mes := &message.Message{
		Level:     "success",
		AlertName: "",
		Text:      "",
		Image:     "",
	}

	err := a.Send(mes)
	require.Error(t, err)
	assert.Equal(t, "error send request, err1", err.Error())
}

func TestSend_error_status_code(t *testing.T) {
	m := &webHookCoreMock{}

	a := &AlertManager{
		whCore: m,
		logger: zap.NewNop(),
	}

	resp := &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer(nil)),
		StatusCode: 0,
	}

	m.On("Send", mock.Anything, mock.Anything).Return(resp, nil)

	mes := &message.Message{
		Level:     "success",
		AlertName: "",
		Text:      "",
		Image:     "",
	}

	err := a.Send(mes)
	require.Error(t, err)
	assert.Equal(t, "unexpected response status code 0", err.Error())
}

func TestSend(t *testing.T) {
	m := &webHookCoreMock{}

	a := &AlertManager{
		whCore: m,
		logger: zap.NewNop(),
	}

	resp := &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer(nil)),
		StatusCode: 200,
	}

	m.On("Send", mock.Anything, mock.Anything).Return(resp, nil)

	mes := &message.Message{
		Level:     "success",
		AlertName: "",
		Text:      "",
		Image:     "",
		Fields:    map[string]string{"foo": "bar"},
	}

	err := a.Send(mes)
	require.NoError(t, err)
}
