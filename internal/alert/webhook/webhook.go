package webhook

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type WebHook struct {
	logger *zap.Logger
	name   string
	url    string
}

func New(cfg config.ChannelsWebHook, logger *zap.Logger) (*WebHook, error) {
	m := &WebHook{
		logger: logger,
		name:   cfg.Name,
		url:    cfg.URL,
	}

	return m, nil
}

func (m *WebHook) Name() string {
	return m.name
}

func (m *WebHook) Send(message string) error {

	var bodyRaw []byte

	bodyRaw = []byte(fmt.Sprintf("{\"text\":\"%s\"}", message))

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

	m.logger.Debug("response from webhook", zap.String("name", m.name), zap.Int("statuscode", res.StatusCode), zap.ByteString("body", responseBody))

	return nil
}
