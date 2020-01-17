package slack

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type httpClientMock struct {
	mock.Mock
}

func (m *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSlack_Name(t *testing.T) {
	s := &Slack{name: "name1"}
	assert.Equal(t, "name1", s.Name())
}

func TestSlack_send(t *testing.T) {
	message := slackMessage{
		Text: "text1",
		Blocks: []block{
			{
				Type: "b1",
				Section: &blockText{
					Type:  "section",
					Emoji: false,
					Text:  "text1",
				},
				Context: nil,
			},
		},
	}

	messageJson, _ := json.Marshal(message)

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
	}

	clientMock := &httpClientMock{}
	clientMock.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)
		body, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, messageJson, body)

	}).Return(res, nil)

	s := &Slack{
		client: clientMock,
	}

	err := s.send(message)
	require.NoError(t, err)
}

func TestSlack_SendSuccess(t *testing.T) {
	message := slackMessage{
		Text: "successpreifx-alert-text",
		Blocks: []block{
			{
				Type: "section",
				Section: &blockText{
					Type:  "mrkdwn",
					Emoji: false,
					Text:  "successpreifx-alert-text",
				},
				Context: nil,
			},
			{
				Type: "context",
				Context: &[]blockText{
					{
						Type:  "mrkdwn",
						Emoji: false,
						Text:  "Alert name: alert-name",
					},
				},
			},
		},
	}

	messageJson, _ := json.Marshal(message)

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
	}

	clientMock := &httpClientMock{}
	clientMock.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)
		body, _ := ioutil.ReadAll(req.Body)

		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, string(messageJson), string(body))

	}).Return(res, nil)

	s := &Slack{
		messagePrefixSuccess: "successpreifx-",
		messagePrefixError:   "errorpreifx-",
		client:               clientMock,
	}

	err := s.SendSuccess("alert-name", "alert-text")
	require.NoError(t, err)
}

func TestSlack_SendError(t *testing.T) {
	message := slackMessage{
		Text: "errorpreifx-alert-text",
		Blocks: []block{
			{
				Type: "section",
				Section: &blockText{
					Type:  "mrkdwn",
					Emoji: false,
					Text:  "errorpreifx-alert-text",
				},
				Context: nil,
			},
			{
				Type: "context",
				Context: &[]blockText{
					{
						Type:  "mrkdwn",
						Emoji: false,
						Text:  "Alert name: alert-name",
					},
				},
			},
		},
	}

	messageJson, _ := json.Marshal(message)

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
	}

	clientMock := &httpClientMock{}
	clientMock.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)
		body, _ := ioutil.ReadAll(req.Body)

		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, string(messageJson), string(body))

	}).Return(res, nil)

	s := &Slack{
		messagePrefixSuccess: "successpreifx-",
		messagePrefixError:   "errorpreifx-",
		client:               clientMock,
	}

	err := s.SendError("alert-name", "alert-text")
	require.NoError(t, err)
}

func TestSlack_Send(t *testing.T) {
	message := slackMessage{
		Text: "alert-text",
		Blocks: []block{
			{
				Type: "section",
				Section: &blockText{
					Type:  "mrkdwn",
					Emoji: false,
					Text:  "alert-text",
				},
				Context: nil,
			},
			{
				Type: "context",
				Context: &[]blockText{
					{
						Type:  "mrkdwn",
						Emoji: false,
						Text:  "Alert name: alert-name",
					},
				},
			},
		},
	}

	messageJson, _ := json.Marshal(message)

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
	}

	clientMock := &httpClientMock{}
	clientMock.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)
		body, _ := ioutil.ReadAll(req.Body)

		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, string(messageJson), string(body))

	}).Return(res, nil)

	s := &Slack{
		messagePrefixSuccess: "successpreifx-",
		messagePrefixError:   "errorpreifx-",
		client:               clientMock,
	}

	err := s.Send("alert-name", "alert-text")
	require.NoError(t, err)
}
