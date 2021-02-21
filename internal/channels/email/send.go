package email

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/balerter/balerter/internal/message"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Send implements
func (e *Email) Send(mes *message.Message) error {
	var err error
	server := mail.NewSMTPClient()

	server.Host = e.conf.Host
	server.Port, err = strconv.Atoi(e.conf.Port)
	if err != nil {
		return err
	}
	server.Username = e.conf.Username
	server.Password = e.conf.Password
	timeout := e.conf.Timeout
	if timeout < 1 {
		timeout = 10
	}
	server.ConnectTimeout = time.Duration(timeout) * time.Second
	server.SendTimeout = time.Duration(timeout) * time.Second

	switch strings.ToLower(e.conf.Secure) {
	case "none":
		server.Encryption = mail.EncryptionNone
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	default:
		server.Encryption = mail.EncryptionTLS
	}

	if server.Port == 465 && e.conf.Secure == "" {
		server.Encryption = mail.EncryptionSSL
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	to := strings.Split(e.conf.To, ";")
	subject := fmt.Sprintf("[%s/%s]", mes.AlertName, mes.Level)
	email.SetFrom(e.conf.From).AddTo(to...).SetSubject(subject)

	if len(e.conf.Cc) > 0 {
		cc := strings.Split(e.conf.Cc, ";")
		email.AddCc(cc...)
	}

	email.SetBody(mail.TextHTML, mes.Text)

	if len(mes.Image) > 0 {
		img := base64.StdEncoding.EncodeToString([]byte(mes.Image))
		email.AddAttachmentBase64(img, mes.AlertName+".png")
	}

	return email.Send(smtpClient)
}
