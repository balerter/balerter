package notify

import (
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/message"
	systemNotify "github.com/martinlindhe/notify"
	"go.uber.org/zap"
)

// Notify represents a channel of type Notify
type Notify struct {
	name        string
	iconSuccess string
	iconError   string
	iconWarning string
	ignore      bool
}

// New creates new Notify channel
func New(cfg notify.Notify, _ *zap.Logger) (*Notify, error) {
	p := &Notify{
		name:        cfg.Name,
		iconSuccess: cfg.Icons.Success,
		iconError:   cfg.Icons.Error,
		iconWarning: cfg.Icons.Warning,
		ignore:      cfg.Ignore,
	}

	return p, nil
}

// Name returns the channel name
func (p *Notify) Name() string {
	return p.name
}

// Send message to the channel
func (p *Notify) Send(mes *message.Message) error {
	systemNotify.Notify("Balerter", mes.Level, mes.Text, p.getIconByLevel(mes.Level))

	return nil
}

func (p *Notify) getIconByLevel(l string) string {
	switch l {
	case "success":
		return p.iconSuccess
	case "warning":
		return p.iconWarning
	case "error":
		return p.iconError
	}

	return ""
}

func (p *Notify) Ignore() bool {
	return p.ignore
}
