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

type slackMessage struct {
	Text string `json:"text"`
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

func (m *Slack) Send(message string) error {

	mes := slackMessage{
		Text: message,
	}

	bodyRaw, err := json.Marshal(mes)
	if err != nil {
		return err
	}

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
