package slack

import (
	"bytes"
	"encoding/json"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Slack struct {
	logger               *zap.Logger
	name                 string
	url                  string
	messagePrefixSuccess string
	messagePrefixError   string
	client               httpClient
}

func New(cfg config.ChannelSlack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:               logger,
		name:                 cfg.Name,
		url:                  cfg.URL,
		messagePrefixSuccess: cfg.MessagePrefixSuccess,
		messagePrefixError:   cfg.MessagePrefixError,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	return m, nil
}

func (m *Slack) Name() string {
	return m.name
}

func (m *Slack) SendSuccess(name, message string) error {
	mes := createSlackMessage(name, m.messagePrefixSuccess+message)

	return m.send(mes)
}

func (m *Slack) SendError(name, message string) error {
	mes := createSlackMessage(name, m.messagePrefixError+message)

	return m.send(mes)
}

func (m *Slack) Send(name, message string) error {
	mes := createSlackMessage(name, message)

	return m.send(mes)
}

func (m *Slack) send(message slackMessage) error {
	bodyRaw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	body := bytes.NewReader(bodyRaw)
	req, err := http.NewRequest(http.MethodPost, m.url, body)
	if err != nil {
		return err
	}

	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
