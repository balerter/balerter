package notify

import (
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/message"
	systemNotify "github.com/martinlindhe/notify"
	"go.uber.org/zap"
)

type Notify struct {
	name        string
	iconSuccess string
	iconError   string
	iconWarning string
}

func New(cfg *notify.Notify, _ *zap.Logger) (*Notify, error) {
	p := &Notify{
		name:        cfg.Name,
		iconSuccess: cfg.Icons.Success,
		iconError:   cfg.Icons.Error,
		iconWarning: cfg.Icons.Warning,
	}

	return p, nil
}

func (p *Notify) Name() string {
	return p.name
}

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
