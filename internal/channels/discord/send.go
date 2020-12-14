package discord

import (
	"github.com/balerter/balerter/internal/message"
)

// Send implements
func (d *Discord) Send(mes *message.Message) error {
	_, err := d.session.SendMessage(d.chanID, mes.Text, nil)
	if err != nil {
		return err
	}
	return nil
}
