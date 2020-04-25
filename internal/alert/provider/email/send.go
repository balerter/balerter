package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (e *Email) Send(message *message.Message) error {
	var (
		c    *smtp.Client
		conn net.Conn
		err  error
	)
	if e.conf.ServerPort == "465" {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}

		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = e.conf.ServerName
		}

		conn, err = tls.Dial("tcp", e.conf.ServerName, tlsConfig)
		if err != nil {
			return errors.Wrap(err, "establish TLS connection to server")
		}
	} else {
		var (
			d   = net.Dialer{}
			err error
		)
		conn, err = d.Dial("tcp", e.conf.ServerName)
		if err != nil {
			return errors.Wrap(err, "establish connection to server")
		}
	}
	c, err = smtp.NewClient(conn, e.conf.ServerName)
	if err != nil {
		conn.Close()
		return errors.Wrap(err, "create SMTP client")
	}
	defer func() {
		if err := c.Quit(); err != nil {
			e.logger.Warn("email client", zap.String("smtp", "failed to close SMTP connection"), zap.Error(err))
		}
	}()

	if ok, mech := c.Extension("AUTH"); ok {
		auth, err := e.auth(mech)
		if err != nil {
			return errors.Wrap(err, "find auth mechanism")
		}
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return errors.Wrapf(err, "%T auth", auth)
			}
		}
	}

	addrs, err := mail.ParseAddressList(e.conf.From)
	if err != nil {
		return errors.Wrap(err, "parse 'from' addresses")
	}
	if len(addrs) != 1 {
		return errors.Errorf("must be exactly one 'from' address (got: %d)", len(addrs))
	}
	if err = c.Mail(addrs[0].Address); err != nil {
		return errors.Wrap(err, "send MAIL command")
	}
	addrs, err = mail.ParseAddressList(e.conf.To)
	if err != nil {
		return errors.Wrapf(err, "parse 'to' addresses")
	}
	for _, addr := range addrs {
		if err = c.Rcpt(addr.Address); err != nil {
			return errors.Wrapf(err, "send RCPT command")
		}
	}

	msg, err := c.Data()
	if err != nil {
		return errors.Wrapf(err, "send DATA command")
	}
	defer msg.Close()

	buffer := &bytes.Buffer{}
	multipartBuffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(multipartBuffer)

	fmt.Fprintf(buffer, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(buffer, "Content-Type: multipart/alternative;  boundary=%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(buffer, "MIME-Version: 1.0\r\n\r\n")

	_, err = msg.Write(buffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "write headers")
	}

	if len(message.Text) > 0 {
		w, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-Transfer-Encoding": {"quoted-printable"},
			"Content-Type":              {"text/plain; charset=UTF-8"},
		})
		if err != nil {
			return errors.Wrap(err, "create part for text template")
		}

		qw := quotedprintable.NewWriter(w)
		_, err = qw.Write([]byte(message.Text))
		if err != nil {
			return errors.Wrap(err, "write text part")
		}
		err = qw.Close()
		if err != nil {
			return errors.Wrap(err, "close text part")
		}
	}

	err = multipartWriter.Close()
	if err != nil {
		return errors.Wrap(err, "close multipartWriter")
	}

	_, err = msg.Write(multipartBuffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "write body buffer")
	}

	return nil
}
