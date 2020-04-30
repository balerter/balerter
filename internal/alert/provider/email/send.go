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
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Encryption = mail.EncryptionTLS

	if e.conf.Port == "465" {
		server.Encryption = mail.EncryptionSSL
	}
	if e.conf.WithoutTLS {
		server.Encryption = mail.EncryptionNone
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	subject := fmt.Sprintf("[%s/%s] %s\r\n", message.AlertName, message.Level, strings.Join(message.Fields, ","))
	email.SetFrom(e.conf.From).
		AddTo(e.conf.To).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, message.Text)

	if len(message.Image) > 0 {
		img := base64.StdEncoding.EncodeToString([]byte(message.Image))
		email.AddAttachmentBase64(img, message.AlertName+".png")
	}

	return email.Send(smtpClient)
}
