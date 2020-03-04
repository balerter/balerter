package syslog

import (
	"encoding/json"
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
)

func (sl *Syslog) Send(mes *message.Message) error {
	data, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling message, %w", err)
	}

	n, err := sl.w.Write(data)
	if err != nil {
		return fmt.Errorf("error write message to syslog, %w", err)
	}

	if n != len(data) {
		return fmt.Errorf("write unexpected bytes count to syslog: %d, expect %d", n, len(data))
	}

	return nil
}
