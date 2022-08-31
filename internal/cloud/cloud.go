package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/config/cloud"
	"github.com/balerter/balerter/internal/message"

	"go.uber.org/zap"
)

const (
	defaultTimeout = time.Second * 5
)

var (
	baseURL = "https://api.balerter.com/s"
	logger  = zap.NewNop()
	version = ""
	enabled = false
	token   = ""
	client  *http.Client
)

func Init(cfg *cloud.Cloud, v string, l *zap.Logger) {
	logger = l.With(zap.String("service", "cloud"))

	token = cfg.Token
	version = v
	enabled = true

	if cfg.Server != "" {
		baseURL = cfg.Server
	}

	tm := defaultTimeout
	if cfg.Timeout > 0 {
		tm = time.Millisecond * time.Duration(cfg.Timeout)
	}
	client = &http.Client{
		Timeout: tm,
	}
}

func SendStart() {
	if !enabled {
		return
	}

	err := rawRequest("/start", http.NoBody)
	if err != nil {
		logger.Error("error send start event", zap.Error(err))
	}
}
func SendStop() {
	if !enabled {
		return
	}

	err := rawRequest("/stop", http.NoBody)
	if err != nil {
		logger.Error("error send stop event", zap.Error(err))
	}
}

func SendAlert(alertName, group string, level alert.Level, levelWasUpdated bool) {
	if !enabled {
		return
	}

	d := struct {
		Name            string      `json:"name,omitempty"`
		Group           string      `json:"group,omitempty"`
		Level           alert.Level `json:"level,omitempty"`
		LevelWasUpdated bool        `json:"level_was_updated,omitempty"`
	}{
		Name:            alertName,
		Group:           group,
		Level:           level,
		LevelWasUpdated: levelWasUpdated,
	}

	alertData, errMarshal := json.Marshal(d)
	if errMarshal != nil {
		logger.Error("error marshal alert", zap.Error(errMarshal))
		return
	}

	err := rawRequest("/alert", bytes.NewReader(alertData))
	if err != nil {
		logger.Error("error send alert event", zap.Error(err))
	}
}

func SendMessage(mes *message.Message) error {
	if !enabled {
		return fmt.Errorf("cloud is not enabled")
	}

	data, errMarshal := json.Marshal(mes)
	if errMarshal != nil {
		return fmt.Errorf("error marshal message, %w", errMarshal)
	}

	err := rawRequest("/message", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error send alert event, %w", err)
	}
	return nil
}

func rawRequest(method string, body io.Reader) error {
	req, errReq := http.NewRequest(http.MethodPost, baseURL+method, body)
	if errReq != nil {
		return fmt.Errorf("error create request: %w", errReq)
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("User-Agent", "balerter-"+version)

	resp, errDo := client.Do(req)
	if errDo != nil {
		return fmt.Errorf("error do request: %w", errDo)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response status code: %d", resp.StatusCode)
	}

	return nil
}
