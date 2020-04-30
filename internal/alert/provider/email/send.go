package email

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/balerter/balerter/internal/alert/message"
	mail "github.com/xhit/go-simple-mail/v2"
)

func (e *Email) Send(message *message.Message) error {
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
	subject := fmt.Sprintf("[%s/%s] %s\r\n", message.AlertName, message.Level, strings.Join(message.Fields, ","))
	email.SetFrom(e.conf.From).AddTo(e.conf.To).SetSubject(subject)

	if len(e.conf.Cc) > 0 {
		email.AddCc(e.conf.Cc)
	}

	email.SetBody(mail.TextHTML, message.Text)

	if len(message.Image) > 0 {
		img := base64.StdEncoding.EncodeToString([]byte(message.Image))
		email.AddAttachmentBase64(img, message.AlertName+".png")
	}

	return email.Send(smtpClient)
}
