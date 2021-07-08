package twiliovoice

import (
	"github.com/balerter/balerter/internal/config/channels/twiliovoice"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//go:generate moq -out http_client_mock.go -skip-ensure -fmt goimports . httpClient

const (
	apiPrefix            = "https://api.twilio.com/2010-04-01"
	defaultClientTimeout = time.Second * 30
)

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type TwilioVoice struct {
	name   string
	sid    string
	token  string
	from   string
	to     string
	twiML  string
	ignore bool

	apiPrefix string
	client    httpClient
	timeout   time.Duration

	logger *zap.Logger
}

func New(cfg twiliovoice.Twilio, logger *zap.Logger) (*TwilioVoice, error) {
	tw := &TwilioVoice{
		name:      cfg.Name,
		sid:       cfg.SID,
		token:     cfg.Token,
		from:      cfg.From,
		to:        cfg.To,
		twiML:     cfg.TwiML,
		ignore:    cfg.Ignore,
		timeout:   time.Millisecond * time.Duration(cfg.Timeout),
		logger:    logger,
		apiPrefix: apiPrefix,
	}

	if tw.timeout == 0 {
		tw.timeout = defaultClientTimeout
	}

	tw.client = &http.Client{}

	return tw, nil
}

func (tw *TwilioVoice) Name() string {
	return tw.name
}

func (tw *TwilioVoice) Ignore() bool {
	return tw.ignore
}
