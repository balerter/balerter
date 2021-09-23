package discord

import (
	"fmt"
	"github.com/balerter/balerter/internal/message"
)

// Send implements
func (d *Discord) Send(mes *message.Message) error {
	if len(mes.Fields) > 0 {
		mes.Text += "\n\n"
		for k, v := range mes.Fields {
			mes.Text += fmt.Sprintf("%s = %s\n", k, v)
		}
	}

	_, err := d.session.SendMessage(d.chanID, mes.Text, nil)
	if err != nil {
		return err
	}
	return nil
}
