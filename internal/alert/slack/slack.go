package slack

import (
	"bytes"
	"encoding/json"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Slack struct {
	logger *zap.Logger
	name   string
	url    string
}

func New(cfg config.ChannelSlack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger: logger,
		name:   cfg.Name,
		url:    cfg.URL,
	}

	return m, nil
}

func (m *Slack) Name() string {
	return m.name
}

func (m *Slack) SendSuccess(name, message string) error {
	mes := createSlackMessage(name, ":eight_spoked_asterisk: "+message)

	return m.send(mes)
}

func (m *Slack) SendError(name, message string) error {
	mes := createSlackMessage(name, ":sos: "+message)

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

	m.logger.Debug("slack message", zap.String("message", string(bodyRaw)))

	body := bytes.NewReader(bodyRaw)

	res, err := http.Post(m.url, "application/json", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	m.logger.Debug("response from slack", zap.String("name", m.name), zap.Int("statuscode", res.StatusCode), zap.ByteString("body", responseBody))

	return nil
}
